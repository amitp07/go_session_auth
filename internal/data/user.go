package data

import (
	"session-auth/internal/dto"
	"session-auth/internal/utils"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username  string       `json:"username" gorm:"column:user_name;unique;not null;"`
	Password  string       `json:"password" gorm:"not null"`
	Groups    []UserGroups `json:"groups" gorm:"many2many:user_group_members"`
	CreatedAt time.Time    `json:"created_at" gorm:"default:NOW();"`
	UpdatedAt time.Time    `json:"updated_at" gorm:"default:NOW();"`
}

func (d *Data) CreateUser(u dto.UserRequest) (uuid.UUID, error) {
	password, err := utils.HashPassword(u.Password)

	if err != nil {
		return uuid.Nil, err
	}

	user := User{
		Username: u.Username,
		Password: password,
	}
	res := d.db.Create(&user)
	if res.Error != nil {
		return uuid.Nil, res.Error
	}

	return user.ID, nil

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
