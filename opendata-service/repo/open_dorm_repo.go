package repo

import (
	"context"
	"opendata-service/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type OpenDormRepository struct {
	collection *mongo.Collection
}

// NewDormRepository kreira novi repo koji samo čita Dorm podatke
func NewDormRepository(db *mongo.Database) *OpenDormRepository {
	return &OpenDormRepository{
		collection: db.Collection("dorms"),
	}
}

// FindAll vraća sve dormove iz baze
func (r *OpenDormRepository) FindAll() ([]model.OpenDorm, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var dorms []model.OpenDorm
	if err := cursor.All(ctx, &dorms); err != nil {
		return nil, err
	}

	return dorms, nil
}

// FindByID vraća dorm po ID-u
func (r *OpenDormRepository) FindByID(id string) (*model.OpenDorm, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var dorm model.OpenDorm
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&dorm)
	if err != nil {
		return nil, err
	}

	return &dorm, nil
}

// Optional: funkcija za filtriranje po kapacitetu ili tipu
func (r *OpenDormRepository) FilterByCapacityOrType(minCapacity, maxCapacity int, dormType string) ([]model.OpenDorm, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	if minCapacity > 0 || maxCapacity > 0 {
		filter["capacity"] = bson.M{}
		if minCapacity > 0 {
			filter["capacity"].(bson.M)["$gte"] = minCapacity
		}
		if maxCapacity > 0 {
			filter["capacity"].(bson.M)["$lte"] = maxCapacity
		}
	}
	if dormType != "" {
		filter["type"] = dormType
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var dorms []model.OpenDorm
	if err := cursor.All(ctx, &dorms); err != nil {
		return nil, err
	}

	return dorms, nil
}
