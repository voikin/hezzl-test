package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/voikin/hezzl-test/internal/domain/event"
	"github.com/voikin/hezzl-test/internal/domain/good"
)

// чтобы не было цикличного импорта из repository
type GoodRepo interface {
	CreateGood(ctx context.Context, name string, projectId int) (good.Good, error)
	UpdateGood(ctx context.Context, name, description string, id, projectId int) (good.Good, error)
	DeleteGood(ctx context.Context, id, projectId int) (good.Good, error)
	GetGood(ctx context.Context, id, projectId int) (good.Good, error)
	GetGoods(ctx context.Context, limit, offset int) ([]good.Good, error)
	UpdateGoodPriority(ctx context.Context, projectID, goodID, newPriority int) ([]good.Good, error)
}

type GoodRepoNats struct {
	GoodRepo
	stream nats.JetStreamContext
}

func NewGoodRepo(repo GoodRepo, js nats.JetStreamContext) *GoodRepoNats {
	cfg := &nats.StreamConfig{
		Name:      "EVENTS",
		Subjects:  []string{"events.>"},
		Retention: nats.WorkQueuePolicy,
	}

	_, err := js.AddStream(cfg)
	if err != nil {
		panic(err)
	}

	return &GoodRepoNats{
		GoodRepo: repo,
		stream:   js,
	}
}

func (grn *GoodRepoNats) UpdateGood(ctx context.Context, name, description string, id, campaignId int) (good.Good, error) {
	good, err := grn.GoodRepo.UpdateGood(ctx, name, description, id, campaignId)

	if err != nil {
		return good, err
	}

	ce := &event.ClickhouseEvent{
		Id:          good.ID,
		ProjectId:   good.ProjectId,
		Name:        good.Name,
		Description: good.Description,
		Priority:    good.Priority,
		Removed:     good.Removed,
		EventTime:   time.Now(),
	}

	grn.sendEvent(ce)

	return good, err
}

func (grn *GoodRepoNats) DeleteGood(ctx context.Context, id, projectId int) (good.Good, error) {
	good, err := grn.GoodRepo.DeleteGood(ctx, id, projectId)

	ce := &event.ClickhouseEvent{
		Id:          good.ID,
		ProjectId:   good.ProjectId,
		Name:        good.Name,
		Description: good.Description,
		Priority:    good.Priority,
		Removed:     good.Removed,
		EventTime:   time.Now(),
	}

	grn.sendEvent(ce)

	return good, err
}

func (grn *GoodRepoNats) UpdateGoodPriority(ctx context.Context, projectID, goodID, newPriority int) ([]good.Good, error) {
	updatedGoods, err := grn.GoodRepo.UpdateGoodPriority(ctx, projectID, goodID, newPriority)
	if err != nil {
		return nil, err
	}

	for _, good := range updatedGoods {
		ce := &event.ClickhouseEvent{
			Id:          good.ID,
			ProjectId:   good.ProjectId,
			Name:        good.Name,
			Description: good.Description,
			Priority:    good.Priority,
			Removed:     good.Removed,
			EventTime:   time.Now(),
		}

		grn.sendEvent(ce)
	}

	return updatedGoods, nil
}

func (grn *GoodRepoNats) sendEvent(event *event.ClickhouseEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		return
	}
	s, err := grn.stream.Publish("events.goods", data)
	fmt.Println(s.Domain, s.Stream, s.Sequence, err, event)
}
