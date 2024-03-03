package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/voikin/hezzl-test/internal/domain/good"
	"github.com/voikin/hezzl-test/internal/utils"
)

type GoodRepo struct {
	db *sql.DB
}

func NewGoodRepo(db *sql.DB) *GoodRepo {
	return &GoodRepo{
		db: db,
	}
}

func (gr *GoodRepo) CreateGood(ctx context.Context, name string, projectID int) (good.Good, error) {
	const fName = "CreateGood"
	tx, err := gr.db.BeginTx(ctx, nil)
	if err != nil {
		return good.Good{}, fmt.Errorf("%s: %w", fName, err)
	}
	defer tx.Rollback()

	var id int
	var createdAt time.Time
	var priority int

	err = tx.QueryRowContext(ctx, "INSERT INTO goods (name, project_id) VALUES ($1, $2) RETURNING id, created_at, priority", name, projectID).Scan(&id, &createdAt, &priority)
	if err != nil {
		return good.Good{}, fmt.Errorf("%s: %w", fName, err)
	}

	err = tx.Commit()
	if err != nil {
		return good.Good{}, fmt.Errorf("%s: %w", fName, err)
	}

	return good.Good{
		ID:        id,
		Name:      name,
		ProjectId: projectID,
		CreatedAt: createdAt,
		Priority:  priority,
	}, nil
}

func (gr *GoodRepo) UpdateGood(ctx context.Context, name, description string, id, projectID int) (good.Good, error) {
	const fName = "UpdateGood"
	tx, err := gr.db.BeginTx(ctx, nil)
	if err != nil {
		return good.Good{}, fmt.Errorf("%s: %w", fName, err)
	}
	defer tx.Rollback()

	var goodFromDB good.Good
	err = tx.QueryRowContext(ctx, "SELECT id, created_at, priority FROM goods WHERE id = $1 AND project_id = $2 FOR UPDATE", id, projectID).Scan(&goodFromDB.ID, &goodFromDB.CreatedAt, &goodFromDB.Priority)
	if err != nil {
		return good.Good{}, utils.ErrGoodNotFound
	}

	_, err = tx.ExecContext(ctx, "UPDATE goods SET name = $1, description = $2 WHERE id = $3 AND project_id = $4", name, description, id, projectID)
	if err != nil {
		return good.Good{}, fmt.Errorf("%s: %w", fName, err)
	}

	err = tx.Commit()
	if err != nil {
		return good.Good{}, fmt.Errorf("%s: %w", fName, err)
	}

	goodFromDB.Name = name
	goodFromDB.Description = description

	return goodFromDB, nil
}

func (gr *GoodRepo) DeleteGood(ctx context.Context, id, projectID int) (good.Good, error) {
	const fName = "DeleteGood"
	tx, err := gr.db.BeginTx(ctx, nil)
	if err != nil {
		return good.Good{}, fmt.Errorf("%s: %w", fName, err)
	}
	defer tx.Rollback()

	var goodFromDB good.Good
	err = tx.QueryRowContext(ctx, "DELETE FROM goods WHERE id = $1 AND project_id = $2 RETURNING id, project_id, name, description, priority, removed, created_at", id, projectID).Scan(&goodFromDB.ID, &goodFromDB.ProjectId, &goodFromDB.Name, &goodFromDB.Description, &goodFromDB.Priority, &goodFromDB.Removed, &goodFromDB.CreatedAt)
	if err != nil {
		return good.Good{}, utils.ErrGoodNotFound
	}

	goodFromDB.Removed = true

	err = tx.Commit()
	if err != nil {
		return good.Good{}, fmt.Errorf("%s: %w", fName, err)
	}

	return goodFromDB, nil
}

func (gr *GoodRepo) GetGood(ctx context.Context, id, projectID int) (good.Good, error) {
	const fName = "GetGood"
	var goodFromDB good.Good

	err := gr.db.QueryRowContext(ctx, "SELECT id, name, description, priority, removed, created_at FROM goods WHERE id = $1 AND project_id = $2", id, projectID).Scan(&goodFromDB.ID, &goodFromDB.Name, &goodFromDB.Description, &goodFromDB.Priority, &goodFromDB.Removed, &goodFromDB.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return good.Good{}, utils.ErrGoodNotFound
		}

		return good.Good{}, fmt.Errorf("%s: %w", fName, err)
	}

	return goodFromDB, nil
}

func (gr *GoodRepo) GetGoods(ctx context.Context, limit, offset int) ([]good.Good, error) {
	const fName = "GetGoods"
	query := fmt.Sprintf("SELECT id, project_id, name, description, priority, removed, created_at FROM goods LIMIT %d OFFSET %d", limit, offset)
	rows, err := gr.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fName, err)
	}
	defer rows.Close()

	var goodsList []good.Good
	for rows.Next() {
		var good good.Good
		err := rows.Scan(&good.ID, &good.ProjectId, &good.Name, &good.Description, &good.Priority, &good.Removed, &good.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fName, err)
		}
		goodsList = append(goodsList, good)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fName, err)
	}

	return goodsList, nil
}

func (gr *GoodRepo) UpdateGoodPriority(ctx context.Context, projectID, goodID, newPriority int) ([]good.Good, error) {
	const fName = "UpdateGoodPriority"

	tx, err := gr.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fName, err)
	}
	defer tx.Rollback()

	var oldPriority int
	err = tx.QueryRowContext(ctx, "SELECT priority FROM goods WHERE id = $1 AND project_id = $2", goodID, projectID).Scan(&oldPriority)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrGoodNotFound
		}

		return nil, fmt.Errorf("%s: %w", fName, err)
	}

	_, err = tx.ExecContext(ctx, "UPDATE goods SET priority = $1 WHERE id = $2 AND project_id = $3", newPriority, goodID, projectID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fName, err)
	}

	_, err = tx.ExecContext(ctx, "UPDATE goods SET priority = priority + 1 WHERE project_id = $1 AND priority >= $2 AND id != $3", projectID, newPriority, goodID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fName, err)
	}

	rows, err := tx.QueryContext(ctx, "SELECT id, project_id, name, description, priority, removed, created_at FROM goods WHERE project_id = $1 ORDER BY priority", projectID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fName, err)
	}
	defer rows.Close()

	var updatedGoods []good.Good
	for rows.Next() {
		var good good.Good
		err := rows.Scan(&good.ID, &good.ProjectId, &good.Name, &good.Description, &good.Priority, &good.Removed, &good.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fName, err)
		}
		updatedGoods = append(updatedGoods, good)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fName, err)
	}

	return updatedGoods, nil
}
