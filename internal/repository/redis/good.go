package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/voikin/hezzl-test/internal/domain/good"
)

// продублировал для избежания цикличного импорта из repository
type GoodRepo interface {
	CreateGood(ctx context.Context, name string, projectId int) (good.Good, error)
	UpdateGood(ctx context.Context, name, description string, id, projectId int) (good.Good, error)
	DeleteGood(ctx context.Context, id, projectId int) (good.Good, error)
	GetGood(ctx context.Context, id, projectId int) (good.Good, error)
	GetGoods(ctx context.Context, limit, offset int) ([]good.Good, error)
	UpdateGoodPriority(ctx context.Context, projectID, goodID, newPriority int) ([]good.Good, error)
}

type RedisGoodRepo struct {
	GoodRepo
	cache *redis.Client
}

func NewRedisGoodRepo(repo GoodRepo, client *redis.Client) *RedisGoodRepo {
	return &RedisGoodRepo{
		GoodRepo: repo,
		cache:    client,
	}
}

func (gr *RedisGoodRepo) CreateGood(ctx context.Context, name string, projectId int) (good.Good, error) {
	return gr.GoodRepo.CreateGood(ctx, name, projectId)
}

func (gr *RedisGoodRepo) GetGoods(ctx context.Context, limit, offset int) ([]good.Good, error) {
	val, err := gr.cache.Get(ctx, "GetGoods").Bytes()
	if err != nil {
		goodsList, err := gr.GoodRepo.GetGoods(ctx, limit, offset)
		if err != nil {
			return nil, err
		}

		data, err := json.Marshal(goodsList)
		if err != nil {
			return goodsList, nil
		}
		gr.cache.SetNX(ctx, "GetGoods", data, _defaultExpiration)
		return goodsList, nil
	}
	goodsList := make([]good.Good, 0)
	err = json.Unmarshal(val, &goodsList)
	if err != nil {
		return nil, err
	}
	return goodsList, err
}

func (gr *RedisGoodRepo) GetGood(ctx context.Context, id, projectId int) (good.Good, error) {
	goodKey := fmt.Sprintf("GetGood-%d-%d", id, projectId)
	val, err := gr.cache.Get(ctx, goodKey).Bytes()

	if err != nil {
		it, err := gr.GoodRepo.GetGood(ctx, id, projectId)
		if err != nil {
			return good.Good{}, err
		}

		data, err := json.Marshal(it)
		if err != nil {
			return it, nil
		}
		gr.cache.SetNX(ctx, goodKey, data, _defaultExpiration)
		return it, nil
	}

	it := good.Good{}
	err = json.Unmarshal(val, &it)
	if err != nil {
		return it, err
	}
	return it, err
}

func (gr *RedisGoodRepo) UpdateGood(ctx context.Context, name, description string, id, projectId int) (good.Good, error) {
	gr.deleteKey(ctx, id, projectId)
	return gr.GoodRepo.UpdateGood(ctx, name, description, id, projectId)
}

func (gr *RedisGoodRepo) DeleteGood(ctx context.Context, id, projectId int) (good.Good, error) {
	gr.deleteKey(ctx, id, projectId)
	return gr.GoodRepo.DeleteGood(ctx, id, projectId)
}

func (gr *RedisGoodRepo) UpdateGoodPriority(ctx context.Context, projectID, goodID, newPriority int) ([]good.Good, error) {
	updatedGoods, err := gr.GoodRepo.UpdateGoodPriority(ctx, projectID, goodID, newPriority)
	if err != nil {
		return nil, err
	}

	for _, good := range updatedGoods {
		gr.deleteKey(ctx, good.ID, projectID)
	}

	return updatedGoods, nil
}

func (gr *RedisGoodRepo) deleteKey(ctx context.Context, id, projectId int) {
	goodKey := fmt.Sprintf("GetGood-%d-%d", id, projectId)
	gr.cache.Del(ctx, goodKey)
	gr.cache.Del(ctx, "GetGoods")
}
