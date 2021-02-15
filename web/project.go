package web

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/grrrben/observe/entity"
	"github.com/grrrben/observe/repo"
	"html/template"
	"net/http"
)

func ServeProject(w http.ResponseWriter, r *http.Request) {

	v := mux.Vars(r)
	hash := v["hash"]

	type PageVars struct {
		Project entity.Project
		QrCode  string
	}

	tmplHelper := NewTemplateHelper()
	tFile := tmplHelper.FilePath("project/view.html")
	tmpl := template.Must(template.ParseFiles(tFile))

	glog.Infof("ServeProject, %s", hash)

	pr := repo.NewProjectRepo()
	p, err := pr.GetByHash(hash)

	if err != nil {
		glog.Errorf("Error in database query: msg %s", err)
		eh := ErrorHandler{Code: http.StatusInternalServerError, Message: err.Error()}
		eh.ServeHTTP(w, r)
		return
	}

	vars := PageVars{Project: p, QrCode: fmt.Sprintf("/static/images/qr/%s.png", p.Hash)}
	tmpl.Execute(w, vars)
}
