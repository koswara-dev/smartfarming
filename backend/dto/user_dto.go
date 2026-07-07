package dto

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=admin operator user"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name" binding:"omitempty,min=2,max=100"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Password *string `json:"password" binding:"omitempty,min=6"`
	Role     *string `json:"role" binding:"omitempty,oneof=admin operator user"`
}
