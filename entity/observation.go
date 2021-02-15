package entity

import (
	"errors"
	"time"
)

type Observation struct {
	Id          int64
	ProjectHash string
	Image       string
	Time        time.Time
}

func (o Observation) GetId() (int64, error) {
	if o.Id == 0 {
		return 0, errors.New("Id of observation not set (not saved?)")
	}
	return o.Id, nil
}

func (o Observation) SetId(id int64) {
	o.Id = id
}
