package entity

import (
	"errors"
	"time"
)

type Observation struct {
	Id          int64
	ProjectHash string
	Image       string
	DateCreated time.Time
}

func (o Observation) GetId() (int64, error) {
	if o.Id == 0 {
		return 0, errors.New("id of observation not set (not saved?)")
	}
	return o.Id, nil
}

func (o Observation) SetId(id int64) {
	o.Id = id
}

func (o Observation) CreatedFormattedDate() string {
	return o.DateCreated.Format("02-01-2006")
}
