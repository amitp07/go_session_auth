package data

import "github.com/google/uuid"

type UserGroups struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string    `json:"name" gorm:"unique;not null"`
	Roles       []Role    `json:"roles" gorm:"many2many:user_groups_roles"`
	Description string    `json:"description"`
}
