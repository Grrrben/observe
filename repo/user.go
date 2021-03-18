package repo

import (
	"errors"
	"fmt"
	"github.com/grrrben/observe/entity"
)

type UserRepository struct {
	*Db
}

func NewUserRepo() *UserRepository {
	return new(UserRepository)
}

func (rep *UserRepository) Get(id int64) (entity.User, error) {
	var u entity.User

	q := `SELECT id, username, email, password, date_created as created FROM "user" WHERE id = $1;`
	c := rep.getConnection()
	defer c.Close()

	err := c.QueryRow(q, id).Scan(&u.Id, &u.Username, &u.Email, &u.Password, &u.DateCreated)
	if err != nil {
		return u, err
	}
	defer c.Close()

	return u, nil
}

func (rep *UserRepository) GetBy(sql string, args ...interface{}) ([]entity.Observation, error) {
	var es []entity.Observation

	return es, errors.New("not implemented")
}

func (rep *UserRepository) Save(e entity.Entity) (entity.Observation, error) {

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
