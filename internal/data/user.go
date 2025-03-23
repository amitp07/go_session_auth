package data

import (
	"session-auth/internal/dto"
	"session-auth/internal/utils"
)

type User struct {
	ID       string `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username string `json:"username" gorm:"column:user_name;not null;"`
	Password string `json:"password" gorm:"not null"`
}

func (d *Data) CreateUser(u dto.UserRequest) (string, error) {
	password, err := utils.HashPassword(u.Password)

	if err != nil {
		return "", err
	}

	user := User{
		Username: u.Username,
		Password: password,
	}
	res := d.db.Create(&user)
	if res.Error != nil {
		return "", res.Error
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
