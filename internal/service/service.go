package service

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/voikin/hezzl-test/internal/domain/event"
	"github.com/voikin/hezzl-test/internal/domain/good"
	"github.com/voikin/hezzl-test/internal/domain/project"
	"github.com/voikin/hezzl-test/internal/repository"
	"github.com/voikin/hezzl-test/internal/service/eventSaver"
	goodService "github.com/voikin/hezzl-test/internal/service/good"
	projectService "github.com/voikin/hezzl-test/internal/service/project"
)

type ProjectService interface {
	CreateProject(ctx context.Context, name string) (project.Project, error)
	UpdateProject(ctx context.Context, name string, id int) (project.Project, error)
	DeleteProject(ctx context.Context, id int) (project.Project, error)
	GetProject(ctx context.Context, id int) (project.Project, error)
	GetProjects(ctx context.Context) ([]project.Project, error)
}

type GoodService interface {
	CreateGood(ctx context.Context, name string, projectId int) (good.Good, error)
	UpdateGood(ctx context.Context, name, description string, id, projectId int) (good.Good, error)
	DeleteGood(ctx context.Context, id, projectId int) (good.Good, error)
	GetGood(ctx context.Context, id, projectId int) (good.Good, error)
	GetGoods(ctx context.Context, limit, offset int) ([]good.Good, error)
	UpdateGoodPriority(ctx context.Context, projectID, goodID, newPriority int) ([]good.GoodPriority, error)
}

type EventSaver interface {
	Start(ctx context.Context)
}

type Event interface {
	CreateEvent(ctx context.Context, event []event.ClickhouseEvent) error
}

type Service struct {
	ProjectService
	GoodService
	EventSaver
}

func NewServices(repo *repository.Repository, js nats.JetStreamContext) *Service {
	return &Service{
		ProjectService: projectService.NewProjectService(repo.ProjectRepo),
		GoodService:    goodService.NewGoodService(repo.GoodRepo),
		EventSaver:     eventSaver.NewEventSaver(repo.EventRepo, js),
	}
}
