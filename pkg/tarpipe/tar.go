package tarpipe

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/JonathanLogan/cfgtar/pkg/jsonschema"
	"github.com/JonathanLogan/cfgtar/pkg/schemareg"
	"github.com/JonathanLogan/cfgtar/pkg/tmpfunc"
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
	var r *tar.Reader
	var w *tar.Writer
	r = tar.NewReader(input)
	if output != nil {
		w = tar.NewWriter(output)
	}
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

		temp := template.New("")
		temp.Option("missingkey=error")
		temp.Funcs(tmpfunc.FuncMap)
		temp, errT := temp.Parse(tempData.String())
		if errT != nil {
			return errT
		}
		buf := new(bytes.Buffer)
		if err := temp.Execute(buf, data); err != nil {
			return err
		}
		if w != nil {
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
	}
	if w != nil {
		if err := w.Close(); err != nil {
			return err
		}
	}
	return nil
}
