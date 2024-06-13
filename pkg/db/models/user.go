package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email,omitempty"`
	FirstName string             `json:"first_name" bson:"first_name,omitempty"`
	LastName  string             `json:"last_name" bson:"last_name,omitempty"`
	CreatedAt int64              `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt int64              `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
