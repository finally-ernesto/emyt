package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model `json:"model"`
	Name       string `json:"name"`
	Email      string `json:"email"`
}
