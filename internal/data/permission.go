package data

import "github.com/google/uuid"

type Permission struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"unique;not null"`
	Description string    `json:"description"`
}
