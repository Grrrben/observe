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
	r.HandleFunc("/observation/{hash}/new", loggerMiddleWare(web.ServeObservationForm)).Methods(http.MethodPost, http.MethodGet)
	r.HandleFunc("/project/new", loggerMiddleWare(web.ServeProjectForm)).Methods(http.MethodPost, http.MethodGet)
	r.HandleFunc("/project/{hash}/viewer", loggerMiddleWare(web.ServeProjectViewer)).Methods(http.MethodGet)
	r.HandleFunc("/project/{hash}/viewer/{timeframe}", loggerMiddleWare(web.ServeProjectViewer)).Methods(http.MethodGet)
	r.HandleFunc("/project/{hash}/", loggerMiddleWare(web.ServeProject)).Methods(http.MethodGet)
	r.HandleFunc("/", loggerMiddleWare(web.ServeHomepage)).Methods(http.MethodGet)

	nf := web.ErrorHandler{Code: http.StatusNotFound, Message: "Unknown route"}
	r.NotFoundHandler = &nf

	mna := web.ErrorHandler{Code: http.StatusMethodNotAllowed, Message: "Method not allowed"}
	r.MethodNotAllowedHandler = &mna

	fs := http.FileServer(http.Dir("/app/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", r)
}

func loggerMiddleWare(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		glog.Infof("[%s] Request url %s", r.Method, r.URL.Path)
		f(w, r)
	}
}
