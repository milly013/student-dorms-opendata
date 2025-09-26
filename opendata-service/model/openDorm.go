package model

type OpenDorm struct {
	ID            string
	Name          string
	City          string
	Capacity      int
	Occupied      int
	OccupancyRate float64
	Type          string
	Amenities     []string
	AverageRating float64
	CommentsCount int
	Tags          []string
}
