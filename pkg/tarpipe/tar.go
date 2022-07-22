package tarpipe

import (
	"archive/tar"
	"bytes"
	"io"
	"text/template"
)

func TarPipe(input io.Reader, output io.Writer, data interface{}) error {
	r := tar.NewReader(input)
	w := tar.NewWriter(output)
	for {
		header, err := r.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		tempData := new(bytes.Buffer)
		if _, err := io.Copy(tempData, r); err != nil {
			return err
		}
		temp, errT := template.New("").Parse(tempData.String())
		if errT != nil {
			return errT
		}
		buf := new(bytes.Buffer)
		if err := temp.Execute(buf, data); err != nil {
			return err
		}
		header.Size = int64(buf.Len())
		if err := w.WriteHeader(header); err != nil {
			return err
		}
		if _, err := io.Copy(w, buf); err != nil {
			return err
		}
		if err := w.Flush(); err != nil {
			return err
		}
	}
	if err := w.Close(); err != nil {
		return err
	}
	return nil
}
