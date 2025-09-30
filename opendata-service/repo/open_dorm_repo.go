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

func (r *OpenDormRepository) FindAll() ([]model.OpenDorm, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Privremeni tip koji uključuje Ratings i Comments
	var rawDorms []struct {
		ID        string   `bson:"id"`
		Name      string   `bson:"name"`
		City      string   `bson:"city"`
		Capacity  int      `bson:"capacity"`
		Occupied  int      `bson:"occupied"`
		Type      string   `bson:"type"`
		Amenities []string `bson:"amenities"`
		Price     float64  `bson:"price"`
		Ratings   []struct {
			UserID string  `bson:"user_id"`
			Score  float64 `bson:"score"`
		} `bson:"ratings,omitempty"`
		Comments []struct{} `bson:"comments,omitempty"`
		Tags     []string   `bson:"tags,omitempty"`
	}

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &rawDorms); err != nil {
		return nil, err
	}

	// Mapiranje u OpenDorm i računanje AverageRating + CommentsCount
	var dorms []model.OpenDorm
	for _, d := range rawDorms {
		var avgRating float64
		if len(d.Ratings) > 0 {
			var sum float64
			for _, r := range d.Ratings {
				sum += r.Score
			}
			avgRating = sum / float64(len(d.Ratings))
		}

		dorms = append(dorms, model.OpenDorm{
			ID:            d.ID,
			Name:          d.Name,
			City:          d.City,
			Capacity:      d.Capacity,
			Occupied:      d.Occupied,
			OccupancyRate: float64(d.Occupied) / float64(d.Capacity) * 100,
			Type:          d.Type,
			Amenities:     d.Amenities,
			AverageRating: avgRating,
			CommentsCount: len(d.Comments),
			Tags:          d.Tags,
			Price:         d.Price,
		})
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

// Funkcija vraća mapu: grad -> mapa amenitija i njihovog broja
func GetFacilitiesSummary(dorms []model.OpenDorm) map[string]map[string]int {
	result := make(map[string]map[string]int)

	for _, dorm := range dorms {
		if _, ok := result[dorm.City]; !ok {
			result[dorm.City] = make(map[string]int)
		}

		for _, amenity := range dorm.Amenities {
			result[dorm.City][amenity]++
		}
	}

	return result
}
