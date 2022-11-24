package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model `json:"model"`
	Username   string `json:"username"`
	password   string `json:"password"`
	Name       string `json:"name"`
	Email      string `json:"email"`
}
