package model

import (
	"time"
)

type Task struct {
	ID    string    `json:"id" bson:"_id"`
	Start time.Time `json:"start" bson:"start"`
	End   time.Time `json:"end" bson:"end"`

	Context string `json:"context" bson:"context"`

	Priority string `json:"priority,omitempty" bson:"priority,omitempty"`

	CompleteDate  time.Time `json:"complete-date,omitempty" bson:"CompleteDate,omitempty"`   // 実際の終了日
	CompleteRatio float64   `json:"complete-ratio,omitempty" bson:"CompleteRatio,omitempty"` // 進捗率

	Parent     string   `json:"parent,omitempty" bson:"parent,omitempty"`
	Children   []string `json:"children,omitempty" bson:"children,omitempty"`
	Dependence []string `json:"dependence,omitempty" bson:"dependence,omitempty"` // タスク間の依存関係

	CreatedAt time.Time `json:"created-at" bson:"createdAt"`
	UpdateAt  time.Time `json:"update-at,omitempty" bson:"updateAt,omitempty"`
}
