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

func (r *OpenDormRepository) GetAveragePricePerCity() (map[string]float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$city"},
			{Key: "averagePrice", Value: bson.D{{Key: "$avg", Value: "$price"}}},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	results := make(map[string]float64)
	for cursor.Next(ctx) {
		var result struct {
			City         string  `bson:"_id"`
			AveragePrice float64 `bson:"averagePrice"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		results[result.City] = result.AveragePrice
	}

	return results, nil
}
func (r *OpenDormRepository) GetFreeSpotsPerCity() (map[string]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// MongoDB aggregation pipeline
	pipeline := mongo.Pipeline{
		{{"$group", bson.D{
			{"_id", "$city"},
			{"totalFreeSpots", bson.D{{"$sum", bson.D{{"$subtract", []interface{}{"$capacity", "$occupied"}}}}}},
		}}},
		{{"$sort", bson.D{{"totalFreeSpots", -1}}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make(map[string]int)
	for cursor.Next(ctx) {
		var row struct {
			City           string `bson:"_id"`
			TotalFreeSpots int    `bson:"totalFreeSpots"`
		}
		if err := cursor.Decode(&row); err != nil {
			return nil, err
		}
		result[row.City] = row.TotalFreeSpots
	}

	return result, nil
}

func (r *OpenDormRepository) GetOccupancyPerCity() (map[string]float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{"$group", bson.D{
			{"_id", "$city"},
			{"totalCapacity", bson.D{{"$sum", "$capacity"}}},
			{"totalOccupied", bson.D{{"$sum", "$occupied"}}},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	results := make(map[string]float64)
	for cursor.Next(ctx) {
		var r struct {
			City          string `bson:"_id"`
			TotalCapacity int    `bson:"totalCapacity"`
			TotalOccupied int    `bson:"totalOccupied"`
		}
		if err := cursor.Decode(&r); err != nil {
			return nil, err
		}
		if r.TotalCapacity > 0 {
			results[r.City] = float64(r.TotalOccupied) / float64(r.TotalCapacity) * 100
		} else {
			results[r.City] = 0
		}
	}

	return results, nil
}
