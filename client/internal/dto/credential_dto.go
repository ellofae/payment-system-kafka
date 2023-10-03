package dto

type CredentialDTO struct {
	ID       int    `json:"id"`
	Email    string `json:"email" validate:"required,email,max=128"`
	Password string `json:"password" validate:"required,max=20"`
}
