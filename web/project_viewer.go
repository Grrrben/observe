package web

import (
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/grrrben/observe/entity"
	"github.com/grrrben/observe/repo"
	"html/template"
	"net/http"
)

const month = "month"

func ServeProjectViewer(w http.ResponseWriter, r *http.Request) {
	var timeFrame string

	v := mux.Vars(r)
	hash := v["hash"]
	timeFrame, ok := v["timeframe"]
	if !ok {
		timeFrame = month
	}

	type PageVars struct {
		Project entity.Project
		// [year][month|week][observations]
		Observations map[int]map[int][]entity.Observation
		TimeFrame    string
		Months       map[int]string
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

	mo := mapObservations(obs, timeFrame)
	vars := PageVars{
		Project:      p,
		TimeFrame:    timeFrame,
		Observations: mo,
		Months:       getTranslatedMonths(),
	}
	tmpl.Execute(w, vars)
}

func mapObservations(observations []entity.Observation, timeFrame string) map[int]map[int][]entity.Observation {
	var mapped = make(map[int]map[int][]entity.Observation, len(observations))
	var tf, y int

	for _, ob := range observations {

		if timeFrame == month {
			tf = int(ob.DateCreated.Month())
			y = ob.DateCreated.Year()
		} else {
			// ISOWeek also fetches the year
			y, tf = ob.DateCreated.ISOWeek()
		}

		inner, ok := mapped[y]
		if !ok {
			inner = make(map[int][]entity.Observation)
			mapped[y] = inner
		}

		mapped[y][tf] = append(mapped[y][tf], ob)
	}
	return mapped
}

func getTranslatedMonths() map[int]string {
	m := make(map[int]string, 12)
	m[1] = "January"
	m[2] = "February"
	m[3] = "March"
	m[4] = "April"
	m[5] = "May"
	m[6] = "June"
	m[7] = "July"
	m[8] = "August"
	m[9] = "September"
	m[10] = "October"
	m[11] = "November"
	m[12] = "December"

	return m
}
