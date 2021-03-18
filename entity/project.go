package entity

import "errors"

type Project struct {
	Id                         int64
	UserId                     int64
	User                       User
	Hash                       string
	Name, Description, Address string
	//Lat, Lng float32
}

func (p Project) GetId() (int64, error) {
	if p.Id == 0 {
		return 0, errors.New("id of project not set (not saved?)")
	}
	return p.Id, nil
}

func (p Project) SetId(id int64) {
	p.Id = id
}
