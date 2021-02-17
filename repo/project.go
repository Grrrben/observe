package repo

import (
	"errors"
	"fmt"
	"github.com/grrrben/observe/entity"
)

type ProjectRepository struct {
	*Db
}

func NewProjectRepo() *ProjectRepository {
	return new(ProjectRepository)
}

func (pr *ProjectRepository) Get(id int64) (entity.Project, error) {
	var p = entity.Project{}
	return p, errors.New("not implemented")
}

func (pr *ProjectRepository) GetByHash(hash string) (entity.Project, error) {
	var p entity.Project

	q := `SELECT id, hash, project_name as name, description, address FROM project WHERE hash = $1 LIMIT 1;`
	c := pr.getConnection()
	defer c.Close()

	err := c.QueryRow(q, hash).Scan(&p.Id, &p.Hash, &p.Name, &p.Description, &p.Address)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (pr *ProjectRepository) All() ([]entity.Project, error) {

	q := `SELECT id, hash, project_name as name, description, address FROM project;`
	c := pr.getConnection()
	defer c.Close()

	rows, err := c.Query(q)
	defer rows.Close()

	var ps []entity.Project

	for rows.Next() {
		var p entity.Project
		err = rows.Scan(&p.Id, &p.Hash, &p.Name, &p.Description, &p.Address)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}

	return ps, err
}

func (pr *ProjectRepository) GetBy(sql string, args ...interface{}) ([]entity.Project, error) {
	var es []entity.Project

	return es, errors.New("not implemented")
}

func (pr *ProjectRepository) Save(e entity.Entity) (entity.Project, error) {

	if p, ok := e.(entity.Project); ok {
		q := `INSERT INTO project (hash, project_name, description, address) VALUES ($1, $2, $3, $4);`
		c := pr.getConnection()
		defer c.Close()

		res, err := c.Exec(q, p.Hash, p.Name, p.Description, p.Address)
		if err != nil {
			return p, err
		}

		lii, _ := res.LastInsertId()
		p.SetId(lii)

		return p, nil
	}

	pp := entity.Project{}
	return pp, fmt.Errorf("project type expected, type %T found", e)
}
