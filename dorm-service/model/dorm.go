package model

type Dorm struct {
	ID          string   `json:"id"` // može biti UUID ili ObjectID
	Name        string   `json:"name"`
	Address     string   `json:"address"`
	City        string   `json:"city"`
	Price       float64  `json:"price"`
	Capacity    int      `json:"capacity"`
	Occupied    int      `json:"occupied"`  // trenutno zauzetih mesta
	Type        string   `json:"type"`      // npr. "muški", "ženski", "mešoviti"
	Amenities   []string `json:"amenities"` // npr. "WiFi", "kantina", "biblioteka"
	Description string   `json:"description,omitempty"`

	// Korisnički aspekti
	Ratings     []Rating  `json:"ratings,omitempty"`
	Comments    []Comment `json:"comments,omitempty"`
	Preferences []string  `json:"preferences,omitempty"`
	Users       []string  `json:"users,omitempty" bson:"users,omitempty"`

	FreeSpots int `bson:"-" json:"free_spots"`
}

type Rating struct {
	UserID string  `json:"user_id" bson:"user_id"`
	Score  float64 `json:"score" bson:"score"`
}

type Comment struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
	Date    string `json:"date"` // ISO 8601 format (npr. "2025-09-12T14:30:00Z")
}
