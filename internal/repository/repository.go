package repository

import (
	"context"
	"database/sql"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/voikin/hezzl-test/internal/domain/event"
	"github.com/voikin/hezzl-test/internal/domain/good"
	"github.com/voikin/hezzl-test/internal/domain/project"
	"github.com/voikin/hezzl-test/internal/repository/clickhouse"
	natsRepo "github.com/voikin/hezzl-test/internal/repository/nats"
	"github.com/voikin/hezzl-test/internal/repository/postgres"
	redisRepo "github.com/voikin/hezzl-test/internal/repository/redis"
)

type ProjectRepo interface {
	CreateProject(ctx context.Context, name string) (project.Project, error)
	UpdateProject(ctx context.Context, name string, id int) (project.Project, error)
	DeleteProject(ctx context.Context, id int) (project.Project, error)
	GetProject(ctx context.Context, id int) (project.Project, error)
	GetProjects(ctx context.Context) ([]project.Project, error)
}

type GoodRepo interface {
	CreateGood(ctx context.Context, name string, projectId int) (good.Good, error)
	UpdateGood(ctx context.Context, name, description string, id, projectId int) (good.Good, error)
	DeleteGood(ctx context.Context, id, projectId int) (good.Good, error)
	GetGood(ctx context.Context, id, projectId int) (good.Good, error)
	GetGoods(ctx context.Context, limit, offset int) ([]good.Good, error)
	UpdateGoodPriority(ctx context.Context, projectID, goodID, newPriority int) ([]good.Good, error)
}

type EventRepo interface {
	CreateEvent(ctx context.Context, event []event.ClickhouseEvent) error
}

type Repository struct {
	ProjectRepo
	GoodRepo
	EventRepo
}

func NewRepositories(pgdb *sql.DB, clickhouseConn driver.Conn, js nats.JetStreamContext, client *redis.Client) *Repository {
	pgProjectRepo := postgres.NewProjectRepo(pgdb)
	pgGoodRepo := postgres.NewGoodRepo(pgdb)

	natsPgGoodRepo := natsRepo.NewGoodRepo(pgGoodRepo, js)
	eventRepo := clickhouse.NewEventRepo(clickhouseConn)

	redisPgProjectRepo := redisRepo.NewProjectRepo(pgProjectRepo, client)
	redisNatsPgGoodRepo := redisRepo.NewRedisGoodRepo(natsPgGoodRepo, client)

	return &Repository{
		ProjectRepo: redisPgProjectRepo,
		GoodRepo:    redisNatsPgGoodRepo,
		EventRepo:   eventRepo,
	}
}
