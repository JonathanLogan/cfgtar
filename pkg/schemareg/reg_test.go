package schemareg

import (
	"strings"
	"testing"
)

func TestReg(t *testing.T) {
	reg := New("default")
	if reg.Get([]string{"nothing", "else"}) != "default" {
		t.Error("No default returned")
	}
	td := []string{"nothing", "else"}
	reg.Add(td, strings.Join(td, "."))
	if reg.Get(td) != "nothing.else" {
		t.Error("Precise path not returned")
	}
	if reg.Get([]string{"nothing"}) != "default" {
		t.Error("Should return default")
	}
	td = []string{"nothing", "else", "together"}
	if reg.Get(td) != "nothing.else" {
		t.Error("Sub-path not returned")
	}
	if reg.Get(nil) != "default" {
		t.Error("Default not returned")
	}
}
