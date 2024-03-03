package good

import "time"

type Good struct {
	ID          int       `json:"id" db:"id"`
	ProjectId   int       `json:"project_id" db:"project_id"`
	Priority    int       `json:"priority" db:"priority"`
	Removed     bool      `json:"removed" db:"removed"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type GoodPriority struct {
	ID       int `json:"id" db:"id"`
	Priority int `json:"priority" db:"priority"`
}
