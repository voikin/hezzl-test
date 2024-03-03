package good

import (
	"context"
	"errors"

	"github.com/voikin/hezzl-test/internal/domain/good"
	"github.com/voikin/hezzl-test/internal/repository"
)

type GoodService struct {
	repo repository.GoodRepo
}

func NewGoodService(repo repository.GoodRepo) *GoodService {
	return &GoodService{repo: repo}
}

func (gs *GoodService) CreateGood(ctx context.Context, name string, projectId int) (good.Good, error) {
	if name == "" {
		return good.Good{}, errors.New("the name cannot be empty")
	}
	return gs.repo.CreateGood(ctx, name, projectId)
}

func (gs *GoodService) UpdateGood(ctx context.Context, name, description string, id, projectId int) (good.Good, error) {
	if name == "" {
		return good.Good{}, errors.New("the name cannot be empty")
	}
	return gs.repo.UpdateGood(ctx, name, description, id, projectId)
}

func (gs *GoodService) DeleteGood(ctx context.Context, id, projectId int) (good.Good, error) {
	return gs.repo.DeleteGood(ctx, id, projectId)
}

func (gs *GoodService) GetGood(ctx context.Context, id, projectId int) (good.Good, error) {
	return gs.repo.GetGood(ctx, id, projectId)
}

func (gs *GoodService) GetGoods(ctx context.Context, limit, offset int) ([]good.Good, error) {
	return gs.repo.GetGoods(ctx, limit, offset)
}

func (gs *GoodService) UpdateGoodPriority(ctx context.Context, projectID, goodID, newPriority int) ([]good.GoodPriority, error) {
	updatedGoods, err := gs.repo.UpdateGoodPriority(ctx, projectID, goodID, newPriority)
	if err != nil {
		return nil, err
	}

	updatedPriorities := make([]good.GoodPriority, 0, len(updatedGoods))

	for _, updatedGood := range updatedGoods {
		updatedPriorities = append(updatedPriorities, good.GoodPriority{
			ID:       updatedGood.ID,
			Priority: updatedGood.Priority,
		})
	}

	return updatedPriorities, nil
}
