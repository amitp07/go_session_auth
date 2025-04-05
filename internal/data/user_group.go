package data

import "github.com/google/uuid"

type UserGroup struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string    `json:"name" gorm:"unique;not null"`
	Description string    `json:"description"`
	Roles       []Role    `json:"roles" gorm:"many2many:user_groups_roles"`
}
