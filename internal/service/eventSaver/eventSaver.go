package eventSaver

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/voikin/hezzl-test/internal/domain/event"
	"github.com/voikin/hezzl-test/internal/repository"
)

type EventSaver struct {
	repo repository.EventRepo
	js   nats.JetStreamContext
}

func (es *EventSaver) CheckBatch(ctx context.Context) {
	sub, _ := es.js.PullSubscribe("events.goods",
		"worker",
		nats.PullMaxWaiting(128),
		nats.BindStream("EVENTS"),
	)
	defer sub.Unsubscribe()

	for {
		if _, ok := ctx.Deadline(); ok {
			break
		}

		time.Sleep(time.Second * 10)

		msgs, err := sub.FetchBatch(100, nats.Context(ctx))
		if err != nil {
			continue
		}

		batch := make([]event.ClickhouseEvent, 0, 100)

		for msg := range msgs.Messages() {
			ev := event.ClickhouseEvent{}
			_ = json.Unmarshal(msg.Data, &ev)
			_ = msg.Ack()
			batch = append(batch, ev)
		}

		if len(batch) != 0 {
			err = es.repo.CreateEvent(ctx, batch)
		}

		fmt.Println(batch, err)
	}
}


func (es *EventSaver) Start(ctx context.Context) {
	go es.CheckBatch(ctx)
}

func NewEventSaver(repo repository.EventRepo, js nats.JetStreamContext) *EventSaver {
	es := &EventSaver{
		repo: repo,
		js:   js,
	}
	return es
}
