package event

import (
	"time"
)

type ClickhouseEvent struct {
	Id          int       `json:"Id"`
	ProjectId   int       `json:"ProjectId"`
	Name        string    `json:"Name"`
	Description string    `json:"Description,omitempty"`
	Priority    int       `json:"Priority,omitempty"`
	Removed     bool      `json:"Removed,omitempty"`
	EventTime   time.Time `json:"EventTime"`
}
