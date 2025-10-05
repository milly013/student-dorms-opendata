package model

import "time"

type MoveInRequest struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	StudentID   string    `json:"student_id" bson:"student_id"`
	DormID      string    `json:"dorm_id" bson:"dorm_id"`
	RoomType    string    `json:"room_type" bson:"room_type"`
	Status      string    `json:"status" bson:"status"` // "pending", "approved", "rejected"
	RequestType string    `json:"request_type" bson:"request_type"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}
