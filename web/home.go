package web

import (
	"github.com/golang/glog"
	"github.com/grrrben/observe/entity"
	"github.com/grrrben/observe/repo"
	"html/template"
	"net/http"
)

func ServeHomepage(w http.ResponseWriter, r *http.Request) {

	type PageVars struct {
		Projects []entity.Project
	}

	tmplHelper := NewTemplateHelper()
	tf := tmplHelper.GetExtendedTemplateFiles("home.html")
	tmpl, err := template.ParseFiles(tf...)
	if err != nil {
		glog.Errorf("could not create template file: msg %s", err)
		eh := ErrorHandler{Code: http.StatusInternalServerError, Message: err.Error()}
		eh.ServeHTTP(w, r)
		return
	}

	pr := repo.NewProjectRepo()
	ps, err := pr.All()
	if err != nil {
		glog.Errorf("Error in database query: msg %s", err)
		eh := ErrorHandler{Code: http.StatusInternalServerError, Message: err.Error()}
		eh.ServeHTTP(w, r)
		return
	}

	vars := PageVars{Projects: ps}
	tmpl.Execute(w, vars)
}
