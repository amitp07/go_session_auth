package data

import (
	"fmt"
	"session-auth/internal/dto"
	"session-auth/internal/utils"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID   `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username  string      `json:"username" gorm:"unique;not null;"`
	Password  string      `json:"password" gorm:"not null"`
	Groups    []UserGroup `json:"groups" gorm:"many2many:user_group_members"`
	Roles     []Role      `json:"roles" gorm:"many2many:users_roles"`
	CreatedAt time.Time   `json:"created_at" gorm:"default:NOW();"`
	UpdatedAt time.Time   `json:"updated_at" gorm:"default:NOW();"`
}

func (d *Data) CreateUserWithGroup(u dto.UserRequest, groupName string) error {
	tx := d.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	password, err := utils.HashPassword(u.Password)

	if err != nil {
		return err
	}

	var group UserGroup

	if groupName == "" {
		groupName = "reader"
	}
	if err := tx.Where(UserGroup{Name: groupName}).First(&group).Error; err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Create(&User{
		Username: u.Username,
		Password: password,
		Groups:   []UserGroup{group},
	}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil

}

func (d *Data) CreateUserWithRole(u dto.UserRequest) error {

	tx := d.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	password, err := utils.HashPassword(u.Password)

	if err != nil {
		tx.Rollback()
		return err
	}

	var role Role
	if err := tx.Where(Role{Name: "reader"}).First(&role).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("role error %w", err)
	}

	user := User{
		Username: u.Username,
		Password: password,
		Roles:    []Role{role},
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("create error %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil

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
