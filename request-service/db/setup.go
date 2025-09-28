package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EnsureRequestCollection kreira kolekciju "movein_requests" sa JSON šemom validacije
func EnsureRequestCollection(db *mongo.Database) error {
	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"student_id", "dorm_id", "status", "created_at"},
			"properties": bson.M{
				"student_id": bson.M{"bsonType": "string"},
				"dorm_id":    bson.M{"bsonType": "string"},
				"status":     bson.M{"bsonType": "string"}, // pending, approved, rejected
				"created_at": bson.M{"bsonType": "date"},
			},
		},
	}

	opts := options.CreateCollection().SetValidator(validator)

	// Proveri da li kolekcija već postoji
	collections, err := db.ListCollectionNames(context.Background(), bson.M{"name": "movein_requests"})
	if err != nil {
		return err
	}
	if len(collections) == 0 {
		return db.CreateCollection(context.Background(), "movein_requests", opts)
	}

	return nil
}

// Optional: helper funkcija za test insert
func InsertTestRequest(db *mongo.Database) error {
	req := bson.M{
		"student_id": "test_student_1",
		"dorm_id":    "test_dorm_1",
		"status":     "pending",
		"created_at": time.Now(),
	}
	_, err := db.Collection("movein_requests").InsertOne(context.Background(), req)
	return err
}
