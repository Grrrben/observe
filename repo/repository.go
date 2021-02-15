package repo

import "github.com/grrrben/observe/entity"

type Repository interface {
	// get one based on id
	Get(id uint32) (entity.Entity, error)

	// get a list based on a query
	GetBy(sql string, args ...interface{}) ([]entity.Entity, error)

	// save
	Save(e entity.Entity) (entity.Entity, error)
}
