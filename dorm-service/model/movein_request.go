package model

import "time"

type MoveInRequest struct {
    ID        string    `json:"id" bson:"_id,omitempty"`
    StudentID string    `json:"student_id" bson:"student_id"`
    DormID    string    `json:"dorm_id" bson:"dorm_id"`
    Status    string    `json:"status" bson:"status"` // "pending", "approved", "rejected"
    CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
