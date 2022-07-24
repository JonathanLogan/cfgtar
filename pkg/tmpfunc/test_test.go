package tmpfunc

import (
	"bytes"
	"fmt"
	"testing"
	"text/template"
)

func executeTemplate(templ string, funcMap template.FuncMap, data interface{}) (string, error) {
	tmp := template.New("")
	tmp.Funcs(funcMap)
	tmp, err := tmp.Parse(templ)
	if err != nil {
		return "", err
	}
	tmp.Option("missingkey=error")
	buf := new(bytes.Buffer)
	if err := tmp.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func TestLookup(t *testing.T) {
	templ := "{{ ipv4addr \"+30\" .host }}"
	td := map[string]interface{}{"host": "192.168.0.5/24"}
	out, err := executeTemplate(templ, FuncMap, td)
	if err != nil {
		t.Fatalf("Execute: %s", err)
	}
	fmt.Println(out)
}
