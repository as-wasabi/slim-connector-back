package model

import (
	"time"
)

type Task struct {
	ID      string    `json:"id" bson:"_id"`
	Start   time.Time `json:"start" bson:"start"`
	End     time.Time `json:"end" bson:"end"`
	Context string    `json:"context" bson:"context"`

	Priority string   `json:"priority,omitempty" bson:"priority,omitempty"`
	Parents  []string `json:"parents,omitempty" bson:"parents, omitempty"`

	CreatedAt time.Time `json:"created-at" bson:"createdAt"`
}
