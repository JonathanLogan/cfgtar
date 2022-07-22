package lineconfig

import (
	"bytes"
	"fmt"
	_ "github.com/davecgh/go-spew/spew"
	"testing"
	"text/template"
)

func TestParseLine(t *testing.T) {
	td := "a.bor[1].a (String)=a.b"
	name, validators, content, err := parseLine(td)
	if err != nil {
		t.Errorf("parseLine: %s", err)
	}
	fmt.Printf("'%s' '%s' '%s'\n", name, validators, content)
	fmt.Println(name.ToString())
}

func TestTreeWalk(t *testing.T) {
	tree := make(Tree)
	td := "a.b.c[0].d=value1"
	p, _, value, err := parseLine(td)
	if err != nil {
		t.Fatalf("parseLine: %s", err)
	}
	_, err = tree.walk(p, value)
	if err != nil {
		t.Fatalf("walk set: %s", err)
	}
	e, err := tree.walk(p)
	if err != nil {
		t.Fatalf("walk get: %s", err)
	}
	if e != value {
		t.Error("Values don't match")
	}
}

func TestParser(t *testing.T) {
	td := `a.b (string)=something
		a.a.
        .b (int)=30
        .c (string)=else
		a.c (string)= $a.a.c`
	tree, err := ParseConfig(bytes.NewBuffer([]byte(td)))
	if err != nil {
		t.Fatalf("ParseConfig: %s", err)
	}
	tmp, err := template.New("test").Parse("{{ .a.c }}")
	if err != nil {
		t.Fatalf("Template.Parse: %s", err)
	}
	buf := new(bytes.Buffer)
	if err := tmp.Execute(buf, tree); err != nil {
		t.Fatalf("Exec: %s", err)
	}
	if buf.String() != "else" {
		t.Error("Unexpected output")
	}
	if tree["a"].(Tree)["a"].(Tree)["b"].(int64) != 30 {
		t.Error("Unexpected output 2")
	}
}
