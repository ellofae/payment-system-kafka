package dto

type UserCreationForm struct {
	FirstName string `json:"first_name" validate:"required,max=128" form:"first_name" binding:"required"`
	LastName  string `json:"last_name" validate:"required,max=128" form:"last_name" binding:"required"`
	Email     string `json:"email" validate:"required,email,max=128" form:"email" binding:"required"`
	Password  string `json:"password" validate:"required,min=8,max=20" form:"password" binding:"required"`
}

type UserLoginForm struct {
	Email    string `json:"email" validate:"required,email,max=128" form:"email" binding:"required"`
	Password string `json:"password" validate:"required,max=20" form:"password" binding:"required"`
}
