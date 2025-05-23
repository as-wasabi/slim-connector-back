package model

import (
	"time"
)

type User struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Email     string    `json:"email" bson:"email"`
	CreatedAt time.Time `json:"created-at" bson:"createdAt"`
	UpdateAt  time.Time `json:"update-at" bson:"UpdateAt"`
}
