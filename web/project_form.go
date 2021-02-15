package web

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/grrrben/observe"
	"github.com/grrrben/observe/entity"
	"github.com/grrrben/observe/repo"
	"html/template"
	"net/http"
	"os"
	"strings"
)

func ServeProjectForm(w http.ResponseWriter, r *http.Request) {

	type PageVars struct {
		Message string
	}

	tmplHelper := NewTemplateHelper()
	tFile := tmplHelper.FilePath("project/form.html")

	tmpl := template.Must(template.ParseFiles(tFile))

	if r.Method == http.MethodGet {
		glog.Info("ServeProjectForm, new project")
		e := tmpl.Execute(w, PageVars{Message: ""})
		if e != nil {
			glog.Fatal(e)
		}
		return
	}

	if r.Method == http.MethodPost {
		glog.Info("ServeProjectForm, save project")

		p := entity.Project{
			Name:        r.FormValue("name"),
			Description: r.FormValue("description"),
			Address:     r.FormValue("address"),
		}

		p.Hash = observe.HashString(fmt.Sprintf("%v", p))
		valid, messages := isValid(p)

		if !valid {
			tmpl.Execute(w, PageVars{Message: strings.Join(messages, ", ")})
			return
		}

		pr := repo.NewProjectRepo()
		_, err := pr.Save(p)
		if err != nil {
			tmpl.Execute(w, PageVars{Message: "Unable to save project"})
			return
		}

		go func() {
			// create QR code
			qr := new(observe.Qr)
			urlNewObs := fmt.Sprintf("%s/observation/%s/new", os.Getenv("DOMAIN"), p.Hash)
			qrCode, err := qr.Create(urlNewObs)
			if err != nil {
				glog.Errorf("Could not create QR code %s", err)
				return
			}

			// save
			filename := fmt.Sprintf("./static/images/qr/%s.png", p.Hash)
			err = qr.Save(filename, qrCode)
			if err != nil {
				glog.Errorf("Could not save QR code %s", err)
			}
		}()

		url := fmt.Sprintf("/project/%s/", p.Hash)
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

func isValid(p entity.Project) (ok bool, ers []string) {
	ok = true

	if p.Hash == "" {
		ers = append(ers, "hash missing")
		ok = false
	}
	if p.Name == "" {
		ers = append(ers, "name missing")
		ok = false
	}
	if p.Description == "" {
		ers = append(ers, "description missing")
		ok = false
	}
	if p.Address == "" {
		ers = append(ers, "address missing")
		ok = false
	}

	return
}
