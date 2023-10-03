package dto

type UserCreationForm struct {
	FirstName string `json:"first_name" validate:"required,max=128"`
	LastName  string `json:"last_name" validate:"required,max=128"`
	Email     string `json:"email" validate:"required,email,max=128"`
	Password  string `json:"password" validate:"required,max=20"`
}

type UserLoginForm struct {
	Email    string `json:"email" validate:"required,email,max=128"`
	Password string `json:"password" validate:"required,max=20"`
}
