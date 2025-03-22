package dtos
type UserResponse struct {
	ID uint `json:"id"`
	UserName string `json:"userName"`
	Email string `json:"email"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	DeletedAt string `json:"deletedAt"`
}