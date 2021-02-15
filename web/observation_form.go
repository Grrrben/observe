package web

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/grrrben/observe"
	"github.com/grrrben/observe/entity"
	"github.com/grrrben/observe/repo"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
	"time"
)

func ServeObservationForm(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup

	v := mux.Vars(r)
	hash := v["hash"]

	type PageVars struct {
		Saved   bool
		Message string
		Hash    string
	}
	vars := PageVars{Saved: false, Message: "", Hash: hash}

	tmplHelper := NewTemplateHelper()
	tf := tmplHelper.GetExtendedTemplateFiles("observation/new.html")
	tmpl, err := template.ParseFiles(tf...)
	if err != nil {
		glog.Errorf("could not create template file: msg %s", err)
		eh := ErrorHandler{Code: http.StatusInternalServerError, Message: err.Error()}
		eh.ServeHTTP(w, r)
		return
	}

	if r.Method == http.MethodGet {
		e := tmpl.Execute(w, vars)
		if e != nil {
			glog.Fatal(e)
		}
		return
	}

	if r.Method == http.MethodPost {
		o := entity.Observation{
			ProjectHash: hash,
			DateCreated: time.Now(),
		}

		raw := r.FormValue("image")
		if raw == "" {
			vars.Message = "Required image missing"
			tmpl.Execute(w, vars)
			return
		}

		// try to save it to file
		imageHandler := new(observe.ImageHandler)
		img, err := imageHandler.Base64ToImage(raw)
		if err != nil {
			glog.Error(err)
			eh := ErrorHandler{Code: http.StatusInternalServerError, Message: err.Error()}
			eh.ServeHTTP(w, r)
			return
		}

		p, _ := imageHandler.GetPath(observe.Raw)
		rimg := fmt.Sprintf("%s/%s_%d.png", p, hash, time.Now().Unix())
		err = imageHandler.SavePng(rimg, img)
		if err != nil {
			glog.Error(err)
			eh := ErrorHandler{Code: http.StatusInternalServerError, Message: err.Error()}
			eh.ServeHTTP(w, r)
			return
		}

		// scale the image to smaller sizes
		// TODO which can be done in the background
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := imageHandler.ScaleAllSizes(rimg)
			if err != nil {
				glog.Error(err)
			}
		}()

		o.Image = filepath.Base(rimg)
		or := repo.NewObservationRepo()
		_, err = or.Save(o)
		if err != nil {
			glog.Error(err)
			eh := ErrorHandler{Code: http.StatusInternalServerError, Message: err.Error()}
			eh.ServeHTTP(w, r)
			return
		}

		vars.Saved = true
		wg.Wait()
		tmpl.Execute(w, vars)
	}
}
