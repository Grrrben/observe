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
)

func ServeObservationUploadForm(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	hash := v["hash"]

	type PageVars struct {
		Saved   bool
		Message string
		Hash    string
	}
	vars := PageVars{Saved: false, Message: "", Hash: hash}

	tmplHelper := NewTemplateHelper()
	tf := tmplHelper.GetExtendedTemplateFiles("observation/upload.html")
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

		fp := observe.NewFormParser(r)
		t, err := fp.ParseFieldAsString("datetime")
		if err != nil {
			glog.Error(err)
			vars.Message = err.Error()
			tmpl.Execute(w, vars)
		}

		b, err := fp.GetFileFromField("image")

		imageHandler := new(observe.ImageHandler)
		pngBytes, err := imageHandler.BytesToPng(b)

		p, _ := imageHandler.GetPath(observe.Raw)
		filename := fmt.Sprintf("%s_%d.png", hash, t.Unix())
		rawImgPath := fmt.Sprintf("%s/%s", p, filename)
		err = imageHandler.SavePng(rawImgPath, pngBytes)
		if err != nil {
			glog.Error(err)
			eh := ErrorHandler{Code: http.StatusInternalServerError, Message: err.Error()}
			eh.ServeHTTP(w, r)
			return
		}

		o := entity.Observation{
			ProjectHash: hash,
			Image:       filepath.Base(rawImgPath),
			DateCreated: t,
		}

		// scale the image to smaller sizes
		// TODO which can be done in the background
		go func() {
			image, err := imageHandler.GetImgFromPath(rawImgPath)
			if err != nil {
				glog.Error(err)
				return
			}

			err = imageHandler.ScaleImageAllSizes(image, filename)
			if err != nil {
				glog.Error(err)
				return
			}
		}()

		or := repo.NewObservationRepo()
		_, err = or.Save(o)
		if err != nil {
			glog.Error(err)
			eh := ErrorHandler{Code: http.StatusInternalServerError, Message: "Unable to save the observation"}
			eh.ServeHTTP(w, r)
			return
		}

		vars.Saved = true
		tmpl.Execute(w, vars)
	}
}
