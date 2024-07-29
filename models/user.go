package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name" validate:"required,min=3,max=32"`
	Email    string             `json:"email" bson:"email" validate:"required,email"`
	Password string             `json:"password" bson:"password" validate:"required,min=6"`
}
