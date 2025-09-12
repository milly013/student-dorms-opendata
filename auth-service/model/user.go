package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `json:"username"`
	Email    string             `json:"email"`
	Password string             `json:"password"`
}
