package repo

import (
	"context"
	"dorm-service/model" // ← ZAMENI prema svom modulu
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DormRepository struct {
	collection *mongo.Collection
}

// NewDormRepository vraća novi repozitorijum za dormove
func NewDormRepository(db *mongo.Database) *DormRepository {
	return &DormRepository{
		collection: db.Collection("dorms"),
	}
}

// Insert dodaje novi dom
func (r *DormRepository) Insert(dorm *model.Dorm) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dorm.ID = primitive.NewObjectID().Hex()
	_, err := r.collection.InsertOne(ctx, dorm)
	return err
}

// FindByID vraća dom po ID-ju
func (r *DormRepository) FindByID(id string) (*model.Dorm, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var dorm model.Dorm
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&dorm)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &dorm, nil
}

// FindAll vraća sve domove
func (r *DormRepository) FindAll() ([]model.Dorm, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var dorms []model.Dorm
	if err := cursor.All(ctx, &dorms); err != nil {
		return nil, err
	}
	return dorms, nil
}

// Update ažurira dom po ID-ju
func (r *DormRepository) Update(id string, updated *model.Dorm) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}

	// Kreiramo mapu bez ID polja za update
	updateData := bson.M{
		"name":        updated.Name,
		"address":     updated.Address,
		"city":        updated.City,
		"capacity":    updated.Capacity,
		"occupied":    updated.Occupied,
		"type":        updated.Type,
		"amenities":   updated.Amenities,
		"description": updated.Description,
		"ratings":     updated.Ratings,
		"comments":    updated.Comments,
		"preferences": updated.Preferences,
	}

	update := bson.M{"$set": updateData}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete briše dom po ID-ju
func (r *DormRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}
