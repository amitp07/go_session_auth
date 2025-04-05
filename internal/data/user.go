package data

import (
	"session-auth/internal/dto"
	"session-auth/internal/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username  string       `json:"username" gorm:"column:user_name;unique;not null;"`
	Password  string       `json:"password" gorm:"not null"`
	Groups    []UserGroups `json:"groups" gorm:"many2many:user_group_members"`
	Roles     []Role       `json:"roles" gorm:"many2many:users_roles"`
	CreatedAt time.Time    `json:"created_at" gorm:"default:NOW();"`
	UpdatedAt time.Time    `json:"updated_at" gorm:"default:NOW();"`
}

func (d *Data) CreateUser(u dto.UserRequest) error {

	return d.db.Transaction(func(tx *gorm.DB) error {

		password, err := utils.HashPassword(u.Password)

		if err != nil {
			return err
		}

		var role Role
		if err := tx.Where(Role{Name: "reader"}).First(&role).Error; err != nil {
			return err
		}

		var group UserGroups
		if err := tx.Where(UserGroups{Name: "Reader"}).Scan(&group).Error; err != nil {
			return err
		}

		user := User{
			Username: u.Username,
			Password: password,
			Roles:    []Role{role},
			Groups:   []UserGroups{group},
		}

		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		return nil

	})

}

func (d *Data) GetAllUsers(users *[]User) error {
	err := d.db.Find(users).Error

	if err != nil {
		return err
	}

	return nil
}

func (d *Data) GetUserByUsername(username string) (*User, error) {
	var user User
	res := d.db.Where(User{Username: username}).First(&user)

	if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}
