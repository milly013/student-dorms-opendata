package model

import "time"

type MoveInRequest struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	StudentID   string    `json:"student_id" bson:"student_id"`
	DormID      string    `json:"dorm_id" bson:"dorm_id"`
	RoomType    string    `json:"room_type,omitempty" bson:"room_type,omitempty"`
	Description string    `json:"description,omitempty" bson:"description,omitempty"` // opis kvara
	Status      string    `json:"status" bson:"status"`                               // "pending", "approved", "rejected"
	RequestType string    `json:"request_type" bson:"request_type"`                   // move_in, move_out, issue_report
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}
