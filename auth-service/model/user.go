package model

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"` // Napomena: u pravoj aplikaciji koristi hash
}
