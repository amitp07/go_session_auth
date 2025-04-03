package data

import "github.com/google/uuid"

type Role struct {
	ID          uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string       `json:"name" gorm:"unique;not null"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions" gorm:"many2many:roles_permissions"`
}
