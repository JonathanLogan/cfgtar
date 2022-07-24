package tarpipe

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/JonathanLogan/cfgtar/pkg/jsonschema"
	"github.com/JonathanLogan/cfgtar/pkg/schemareg"
	"io"
	"os"
	"path"
	"strings"
	"text/template"
)

const (
	SchemaFileName = "._config-schema.json"
)

func TarPipe(input io.Reader, output io.Writer, reg *schemareg.Registry) error {
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

		if header.Typeflag == tar.TypeReg && path.Base(header.Name) == SchemaFileName {
			var schema, newData interface{}
			var pErr []string
			var err error
			if err = json.Unmarshal(tempData.Bytes(), &schema); err != nil {
				return err
			}
			if pErr, newData, err = jsonschema.Validate(schema, reg.Get(nil)); err != nil {
				return fmt.Errorf("Validation at '%s': %v %s", header.Name, pErr, err)
			}
			reg.Add(strings.Split(path.Dir(header.Name), string(os.PathSeparator)), newData)
			continue
		}
		data := reg.Get(strings.Split(path.Dir(header.Name), string(os.PathSeparator)))

		temp, errT := template.New("").Parse(tempData.String())
		if errT != nil {
			return errT
		}
		buf := new(bytes.Buffer)
		temp.Option("missingkey=error")
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
