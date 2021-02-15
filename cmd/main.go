package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/grrrben/observe/web"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	flag.Parse()
	glog.Info("App started")

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file; msg %s", err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/observation/{hash}/new", web.ServeObservationForm).Methods(http.MethodPost, http.MethodGet)
	r.HandleFunc("/project/new", web.ServeProjectForm).Methods(http.MethodPost, http.MethodGet)
	r.HandleFunc("/project/{hash}/viewer", web.ServeProjectViewer).Methods(http.MethodGet)
	r.HandleFunc("/project/{hash}/", web.ServeProject).Methods(http.MethodGet)
	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		glog.Info("Requested /")
	})

	nf := web.ErrorHandler{Code: http.StatusNotFound, Message: "Unknown route"}
	r.NotFoundHandler = &nf

	mna := web.ErrorHandler{Code: http.StatusMethodNotAllowed, Message: "Method not allowed"}
	r.MethodNotAllowedHandler = &mna

	fs := http.FileServer(http.Dir("/app/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", r)
}
