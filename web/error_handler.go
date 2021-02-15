package web

import (
	"github.com/golang/glog"
	"html/template"
	"net/http"
)

type ErrorHandler struct {
	Code    uint16
	Message string
}

func (h *ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	glog.Warningf("%d; %s", h.Code, h.Message)

	tmplHelper := NewTemplateHelper()
	files := tmplHelper.GetExtendedTemplateFiles("error.html")

	tmpl, _ := template.ParseFiles(files...)
	tmpl.Execute(w, struct {
		Code    uint16
		Message string
	}{h.Code, h.Message})
}
