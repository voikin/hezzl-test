package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/voikin/hezzl-test/internal/domain/project"
)

// продублировал для избежания цикличного импорта из repository
type ProjectRepo interface {
	CreateProject(ctx context.Context, name string) (project.Project, error)
	UpdateProject(ctx context.Context, name string, id int) (project.Project, error)
	DeleteProject(ctx context.Context, id int) (project.Project, error)
	GetProject(ctx context.Context, id int) (project.Project, error)
	GetProjects(ctx context.Context) ([]project.Project, error)
}

type RedisProjectRepo struct {
	ProjectRepo
	redis *redis.Client
}

func NewProjectRepo(repo ProjectRepo, client *redis.Client) *RedisProjectRepo {
	return &RedisProjectRepo{
		ProjectRepo: repo,
		redis:       client,
	}
}

func (pr *RedisProjectRepo) CreateProject(ctx context.Context, name string) (project.Project, error) {
	return pr.ProjectRepo.CreateProject(ctx, name)
}

func (pr *RedisProjectRepo) GetProjects(ctx context.Context) ([]project.Project, error) {
	redisKey := "GetProjects"
	data, err := pr.redis.Get(ctx, redisKey).Bytes()

	if err != nil {
		projectList, err := pr.ProjectRepo.GetProjects(ctx)
		if err != nil {
			return nil, err
		}

		data, err = json.Marshal(projectList)
		if err != nil {
			return projectList, err
		}

		pr.redis.SetNX(ctx, redisKey, data, _defaultExpiration)
		return projectList, nil
	}

	projectList := make([]project.Project, 0)
	err = json.Unmarshal(data, &projectList)
	if err != nil {
		return nil, err
	}

	return projectList, nil
}

func (pr *RedisProjectRepo) GetProject(ctx context.Context, id int) (project.Project, error) {
	redisKey := fmt.Sprintf("GetProject-%d", id)
	data, err := pr.redis.Get(ctx, redisKey).Bytes()

	if err != nil {
		proj, err := pr.ProjectRepo.GetProject(ctx, id)

		if err != nil {
			return project.Project{}, err
		}

		data, err := json.Marshal(proj)
		if err != nil {
			return project.Project{}, err
		}

		pr.redis.SetNX(ctx, redisKey, data, _defaultExpiration)
		return proj, nil
	}

	proj := project.Project{}
	err = json.Unmarshal(data, &proj)
	if err != nil {
		return project.Project{}, err
	}

	return proj, nil
}

func (pr *RedisProjectRepo) UpdateProject(ctx context.Context, name string, id int) (project.Project, error) {
	pr.deleteKey(ctx, id)
	return pr.ProjectRepo.UpdateProject(ctx, name, id)
}

func (pr *RedisProjectRepo) DeleteProject(ctx context.Context, id int) (project.Project, error) {
	pr.deleteKey(ctx, id)
	return pr.ProjectRepo.DeleteProject(ctx, id)
}

func (pr *RedisProjectRepo) deleteKey(ctx context.Context, id int) {
	redisKey := fmt.Sprintf("GetProject-%d", id)
	pr.redis.Del(ctx, redisKey)
	pr.redis.Del(ctx, "GetProjects")
}
