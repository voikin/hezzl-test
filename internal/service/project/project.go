package project

import (
	"context"
	"errors"

	"github.com/voikin/hezzl-test/internal/domain/project"
	"github.com/voikin/hezzl-test/internal/repository"
)

type ProjectService struct {
	repo repository.ProjectRepo
}

func NewProjectService(repo repository.ProjectRepo) *ProjectService {
	return &ProjectService{repo: repo}
}

func (ps *ProjectService) CreateProject(ctx context.Context, name string) (project.Project, error) {
	if name == "" {
		return project.Project{}, errors.New("the name cannot be empty")
	}
	return ps.repo.CreateProject(ctx, name)
}

func (ps *ProjectService) UpdateProject(ctx context.Context, name string, projectId int) (project.Project, error) {
	if name == "" {
		return project.Project{}, errors.New("the name cannot be empty")
	}
	return ps.repo.UpdateProject(ctx, name, projectId)
}

func (ps *ProjectService) DeleteProject(ctx context.Context, id int) (project.Project, error) {
	return ps.repo.DeleteProject(ctx, id)
}

func (ps *ProjectService) GetProject(ctx context.Context, id int) (project.Project, error) {
	return ps.repo.GetProject(ctx, id)
}

func (ps *ProjectService) GetProjects(ctx context.Context) ([]project.Project, error) {
	return ps.repo.GetProjects(ctx)
}
