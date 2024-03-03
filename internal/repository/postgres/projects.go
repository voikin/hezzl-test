package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/voikin/hezzl-test/internal/domain/project"
	"github.com/voikin/hezzl-test/internal/utils"
)

type ProjectRepo struct {
	db *sql.DB
}

func NewProjectRepo(db *sql.DB) *ProjectRepo {
	return &ProjectRepo{
		db: db,
	}
}

func (pr *ProjectRepo) CreateProject(ctx context.Context, name string) (project.Project, error) {
	const fName = "CreateProject"
	var id int

	err := pr.db.QueryRowContext(ctx, "INSERT INTO projects (name) VALUES ($1) RETURNING id", name).Scan(&id)
	if err != nil {
		return project.Project{}, fmt.Errorf("%s: %w", fName, err)
	}

	return project.Project{ID: id, Name: name}, nil
}

func (pr *ProjectRepo) UpdateProject(ctx context.Context, name string, id int) (project.Project, error) {
	const fName = "UpdateProject"
	tx, err := pr.db.BeginTx(ctx, nil)
	if err != nil {
		return project.Project{}, fmt.Errorf("%s: %w", fName, err)
	}
	defer tx.Rollback()

	var proj project.Project
	err = tx.QueryRowContext(ctx, "SELECT id, name FROM projects WHERE id = $1 FOR UPDATE", id).Scan(&proj.ID, &proj.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return project.Project{}, utils.ErrProjectNotFound
		}
		return project.Project{}, fmt.Errorf("%s: %w", fName, err)
	}

	proj.Name = name
	_, err = tx.ExecContext(ctx, "UPDATE projects SET name = $1 WHERE id = $2", name, id)
	if err != nil {
		return project.Project{}, fmt.Errorf("%s: %w", fName, err)
	}

	err = tx.Commit()
	if err != nil {
		return project.Project{}, fmt.Errorf("%s: %w", fName, err)
	}

	return proj, nil
}

func (pr *ProjectRepo) DeleteProject(ctx context.Context, id int) (project.Project, error) {
	const fName = "DeleteProject"
	tx, err := pr.db.BeginTx(ctx, nil)
	if err != nil {
		return project.Project{}, fmt.Errorf("%s: %w", fName, err)
	}
	defer tx.Rollback()

	var name string
	err = tx.QueryRowContext(ctx, "DELETE FROM projects WHERE id = $1 RETURNING name", id).Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return project.Project{}, utils.ErrProjectNotFound
		}
		return project.Project{}, fmt.Errorf("%s: %w", fName, err)
	}

	err = tx.Commit()
	if err != nil {
		return project.Project{}, fmt.Errorf("%s: %w", fName, err)
	}

	return project.Project{Name: name, ID: id}, nil
}

func (pr *ProjectRepo) GetProject(ctx context.Context, id int) (project.Project, error) {
	const fName = "GetProject"
	var proj project.Project
	err := pr.db.QueryRowContext(ctx, "SELECT id, name FROM projects WHERE id = $1", id).Scan(&proj.ID, &proj.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return project.Project{}, utils.ErrProjectNotFound
		}
		return project.Project{}, fmt.Errorf("%s: %w", fName, err)
	}

	return proj, nil
}

func (pr *ProjectRepo) GetProjects(ctx context.Context) ([]project.Project, error) {
	const fName = "GetProjects"
	rows, err := pr.db.QueryContext(ctx, "SELECT id, name FROM projects")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fName, err)
	}
	defer rows.Close()

	var projects []project.Project
	for rows.Next() {
		var proj project.Project
		err := rows.Scan(&proj.ID, &proj.Name)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fName, err)
		}
		projects = append(projects, proj)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fName, err)
	}

	return projects, nil
}
