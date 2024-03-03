package clickhouse

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/voikin/hezzl-test/internal/domain/event"
)

type EventRepo struct {
	db driver.Conn
}

func NewEventRepo(db driver.Conn) *EventRepo {
	return &EventRepo{db: db}
}

func (er *EventRepo) CreateEvent(ctx context.Context, clickhouseEvents []event.ClickhouseEvent) error {
	insertQuery := "INSERT INTO goods (Id, ProjectId, Name, Description, Priority, Removed, EventTime) VALUES "
	var args []interface{}

	for _, ce := range clickhouseEvents {
		insertQuery += "(?, ?, ?, ?, ?, ?, ?),"
		args = append(args, ce.Id, ce.ProjectId, ce.Name, ce.Description, ce.Priority, ce.Removed, ce.EventTime)
	}

	insertQuery = insertQuery[:len(insertQuery)-1] // чтобы убрать последнюю запятую
	err := er.db.Exec(ctx, insertQuery, args...)
	if err != nil {
		return fmt.Errorf("clickhouse.CreateEvent Exec: %w", err)
	}
	return nil
}
