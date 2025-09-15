package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EnsureDormCollection kreira kolekciju "dorms" sa JSON šemom validacije
func EnsureDormCollection(db *mongo.Database) error {
	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"name", "address", "city", "capacity", "type"},
			"properties": bson.M{
				"name":     bson.M{"bsonType": "string"},
				"address":  bson.M{"bsonType": "string"},
				"city":     bson.M{"bsonType": "string"},
				"capacity": bson.M{"bsonType": "int", "minimum": 0},
				"occupied": bson.M{"bsonType": "int", "minimum": 0},
				"type":     bson.M{"bsonType": "string"}, // možeš dodati enum ako hoćeš da ograničiš vrednosti
				"amenities": bson.M{
					"bsonType": "array",
					"items":    bson.M{"bsonType": "string"},
				},
				"description": bson.M{"bsonType": "string"},
				"preferences": bson.M{
					"bsonType": "array",
					"items":    bson.M{"bsonType": "string"},
				},
				"ratings": bson.M{
					"bsonType": "array",
					"items": bson.M{
						"bsonType": "object",
						"required": []string{"user_id", "score"},
						"properties": bson.M{
							"user_id": bson.M{"bsonType": "string"},
							"score":   bson.M{"bsonType": "double", "minimum": 1.0, "maximum": 5.0},
						},
					},
				},
				"comments": bson.M{
					"bsonType": "array",
					"items": bson.M{
						"bsonType": "object",
						"required": []string{"user_id", "message", "date"},
						"properties": bson.M{
							"user_id": bson.M{"bsonType": "string"},
							"message": bson.M{"bsonType": "string"},
							"date":    bson.M{"bsonType": "string"}, // ISO date kao string
						},
					},
				},
			},
		},
	}

	opts := options.CreateCollection().SetValidator(validator)

	// Proveri da li kolekcija već postoji
	collections, err := db.ListCollectionNames(context.Background(), bson.M{"name": "dorms"})
	if err != nil {
		return err
	}
	if len(collections) == 0 {
		return db.CreateCollection(context.Background(), "dorms", opts)
	}

	return nil
}
