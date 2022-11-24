package models

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"github.com/sethvargo/go-password/password"
)

type User struct {
	gorm.Model `json:"model"`
	Username   string `json:"username"`
	// Password field is not exported
	password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

func (u *User) GeneratePassword() {
	// Generate a password that is 64 characters long with 10 digits, 10 symbols,
	// allowing upper and lower case letters, disallowing repeat characters.
	res, err := password.Generate(64, 10, 10, false, false)
	if err != nil {
		log.Fatal(err)
	}
	u.password = res
}
