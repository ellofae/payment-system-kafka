package entity

import "time"

type Credential struct {
	ID           int       `json:"id" db:"id"`
	Email        string    `json:"email" db:"id"`
	Password     string    `json:"password" db:"password"`
	RegisterDate time.Time `json:"register_date" db:"register_date"`
}
