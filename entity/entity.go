package entity

type Entity interface {
	GetId() (int64, error)
	SetId(id int64)
}
