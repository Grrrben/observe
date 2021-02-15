package web

import (
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/grrrben/observe/entity"
	"github.com/grrrben/observe/repo"
	"html/template"
	"net/http"
)

func ServeProjectViewer(w http.ResponseWriter, r *http.Request) {

	v := mux.Vars(r)
	hash := v["hash"]

	type PageVars struct {
		Project      entity.Project
		Observations []entity.Observation
	}

	tmplHelper := NewTemplateHelper()
	tf := tmplHelper.GetExtendedTemplateFiles("viewer/view.html")
	tmpl, err := template.ParseFiles(tf...)

	pr := repo.NewProjectRepo()
	p, err := pr.GetByHash(hash)
	if err != nil {
		glog.Errorf("Error in database query: msg %s", err)
		eh := ErrorHandler{Code: http.StatusInternalServerError, Message: err.Error()}
		eh.ServeHTTP(w, r)
		return
	}

	or := repo.NewObservationRepo()
	obs, err := or.FindByHash(hash)
	if err != nil {
		glog.Errorf("Error in database query: msg %s", err)
		eh := ErrorHandler{Code: http.StatusInternalServerError, Message: err.Error()}
		eh.ServeHTTP(w, r)
		return
	}

	vars := PageVars{Project: p, Observations: obs}
	tmpl.Execute(w, vars)
}
