package structs

import "time"

type User struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	InsertedAt time.Time `json:"inserted_at"`
}

type UserAuthMethod struct {
	UserID int    `json:"user_id"`
	Type   int    `json:"type"`
	Value  string `json:"value"`
}
