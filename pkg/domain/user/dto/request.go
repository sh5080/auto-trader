package dto

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required" example:"John Doe"`
	Nickname string `json:"nickname" validate:"required" example:"john_doe_nickname"`
	Email    string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"password123"`
}
