package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EnsureUserCollection(db *mongo.Database) error {
	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"username", "email", "password"},
			"properties": bson.M{
				"username": bson.M{"bsonType": "string"},
				"email":    bson.M{"bsonType": "string", "pattern": "^.+@.+$"},
				"password": bson.M{"bsonType": "string", "minLength": 6},
			},
		},
	}

	opts := options.CreateCollection().SetValidator(validator)
	return db.CreateCollection(context.Background(), "users", opts)
}
