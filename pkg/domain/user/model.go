package user

import "time"

type UserExceptPassword struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Nickname  string    `json:"nickname" db:"nickname"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type User struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Nickname  string    `json:"nickname" db:"nickname"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
