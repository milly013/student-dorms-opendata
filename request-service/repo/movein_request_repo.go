package repo

import (
	"context"
	"fmt"
	"request-service/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MoveInRequestRepository struct {
	collection *mongo.Collection
}

// Kreira novi repo za request-service
func NewMoveInRequestRepository(db *mongo.Database) *MoveInRequestRepository {
	return &MoveInRequestRepository{
		collection: db.Collection("movein_requests"),
	}
}

// Create ubacuje novi zahtjev u kolekciju
func (r *MoveInRequestRepository) Create(ctx context.Context, req *model.MoveInRequest) error {
	req.Status = "pending"
	req.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, req)
	return err
}

// GetAll vraća sve zahtjeve
func (r *MoveInRequestRepository) GetAll(ctx context.Context) ([]model.MoveInRequest, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var requests []model.MoveInRequest
	for cursor.Next(ctx) {
		var req model.MoveInRequest
		if err := cursor.Decode(&req); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}

func (r *MoveInRequestRepository) GetByID(ctx context.Context, id string) (*model.MoveInRequest, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var req model.MoveInRequest
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&req)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &req, nil
}

// UpdateStatus mijenja status zahtjeva po ID-u
func (r *MoveInRequestRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	res, err := r.collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// ⭐ Nova funkcija: UpdateStatusWithReason
func (r *MoveInRequestRepository) UpdateStatusWithReason(ctx context.Context, id string, status string, reason string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"status":           status,
			"rejection_reason": reason, // ✅ ispravno ime polja
		},
	}

	res, err := r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// GetByStudent vraća zahtjev jednog studenta
func (r *MoveInRequestRepository) GetByStudent(ctx context.Context, studentID string) (*model.MoveInRequest, error) {
	var req model.MoveInRequest
	err := r.collection.FindOne(ctx, bson.M{"student_id": studentID}).Decode(&req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

// Delete briše zahtjev po ID-u
func (r *MoveInRequestRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ObjectID: %w", err)
	}

	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
