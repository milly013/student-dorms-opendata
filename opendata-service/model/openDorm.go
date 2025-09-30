package model

type OpenDorm struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	City          string   `json:"city"`
	Capacity      int      `json:"capacity"`
	Occupied      int      `json:"occupied"`
	OccupancyRate float64  `json:"occupancy_rate"`
	Type          string   `json:"type"`
	Amenities     []string `json:"amenities"`
	AverageRating float64  `json:"average_rating"`
	CommentsCount int      `json:"comments_count"`
	Tags          []string `json:"tags"`
	Price         float64  `json:"price"`
}
