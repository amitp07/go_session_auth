package database

import (
	"session-auth/internal/data"

	"gorm.io/gorm"
)

func SeedDB(db *gorm.DB) error {
	groups := []data.UserGroup{
		{Name: "admin", Description: "Admin group"},
		{Name: "reader", Description: "read-only group"},
		{Name: "contributor", Description: "Contributor group"},
	}

	roles := []data.Role{
		{Name: "admin"},
		{Name: "super_user"},
		{Name: "contributor"},
		{Name: "reader"},
	}

	permissions := []data.Permission{
		{Name: "read"},
		{Name: "write"},
		{Name: "execute"},
	}

	for _, g := range groups {
		if err := db.FirstOrCreate(&g, data.UserGroup{Name: g.Name}).Error; err != nil {
			return err
		}
	}
	for _, r := range roles {
		if err := db.FirstOrCreate(&r, data.Role{Name: r.Name}).Error; err != nil {
			return err
		}
	}
	for _, p := range permissions {
		if err := db.FirstOrCreate(&p, data.Permission{Name: p.Name}).Error; err != nil {
			return err
		}
	}
	return nil
}
