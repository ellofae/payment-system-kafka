package entity

type User struct {
	ID           int    `json:"id" db:"id"`
	FirstName    string `json:"first_name" db:"first_name"`
	LastName     string `json:"last_name" db:"last_name"`
	CredentialID int    `json:"credential_id" db:"credential_id"`
}
