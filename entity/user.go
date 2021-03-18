package entity

import (
	"errors"
	"time"
)

const (
	role_user = iota
	role_admin
)

type User struct {
	Id, Role                  int64
	Username, Email, Password string
	DateCreated               time.Time
}

func (u User) GetId() (int64, error) {
	if u.Id == 0 {
		return 0, errors.New("id of user not set (not saved?)")
	}
	return u.Id, nil
}

func (u User) SetId(id int64) {
	u.Id = id
}
