package data

import "gorm.io/gorm"

type Data struct {
	db         *gorm.DB
	User       *User
	Permission *Permission
	Role       *Role
	UserGroups *UserGroups
}

func NewModels(db *gorm.DB) *Data {
	return &Data{
		db: db,
	}
}

func (d *Data) MigrateDB() error {
	return d.db.AutoMigrate(&User{})
}
