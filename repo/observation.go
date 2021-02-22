package repo

import (
	"errors"
	"fmt"
	"github.com/grrrben/observe/entity"
)

type ObservationRepository struct {
	*Db
}

func NewObservationRepo() *ObservationRepository {
	return new(ObservationRepository)
}

func (rep *ObservationRepository) Get(id int64) (entity.Observation, error) {
	var p = entity.Observation{}
	return p, errors.New("not implemented")
}

func (rep *ObservationRepository) GetBy(sql string, args ...interface{}) ([]entity.Observation, error) {
	var es []entity.Observation

	return es, errors.New("not implemented")
}

func (rep *ObservationRepository) FindByHash(hash string) ([]entity.Observation, error) {
	var obs []entity.Observation

	q := `SELECT id, project_hash as ProjectHash, image, date_created as created FROM observation WHERE project_hash = $1;`
	c := rep.getConnection()
	defer c.Close()

	rows, err := c.Query(q, hash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var o entity.Observation
		err = rows.Scan(&o.Id, &o.ProjectHash, &o.Image, &o.DateCreated)
		if err != nil {
			return obs, err
		}

		obs = append(obs, o)
	}

	return obs, err
}

func (rep *ObservationRepository) Save(e entity.Entity) (entity.Observation, error) {

	if o, ok := e.(entity.Observation); ok {
		q := `INSERT INTO observation (project_hash, image) VALUES ($1, $2);`
		c := rep.getConnection()
		defer c.Close()

		res, err := c.Exec(q, o.ProjectHash, o.Image)
		if err != nil {
			return o, err
		}

		id, _ := res.LastInsertId()
		o.SetId(id)

		return o, nil
	}
	oo := entity.Observation{}
	return oo, fmt.Errorf("observation type expected, type %T found", e)
}
