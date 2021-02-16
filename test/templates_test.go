package test

import (
	"fmt"
	"github.com/grrrben/observe/web"
	"html/template"
	"testing"
)

type testWriter struct {
	Received []byte
}

func (t testWriter) Write(p []byte) (n int, err error) {
	t.Received = p

	return len(p), nil
}

func TestTemplateInheritance(t *testing.T) {
	TestSetup()
	w := testWriter{}

	tmplHelper := web.NewTemplateHelper()
	files := tmplHelper.GetExtendedTemplateFiles("error.html")

	tmpl, e := template.ParseFiles(files...)
	if e != nil {
		t.Errorf("Received an error during testing; %s", e)
	}

	e = tmpl.Execute(w, struct {
		Code    uint16
		Message string
	}{500, "server error"})

	if e != nil {
		t.Errorf("Received an error during testing; %s", e)
	}

	fmt.Println(string(w.Received))

}
