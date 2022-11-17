package models

import (
	"errors"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint32    `gorm:"primary_key;auto_increment" yaml:"id"`
	Username  string    `gorm:"size:255;not null;unique" yaml:"username"`
	Password  string    `gorm:"size:100;not null;" yaml:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" yaml:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" yaml:"updated_at"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	hasshedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}

	u.Password = string(hasshedPassword)
	return nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "login":
		if u.Password == "" {
			return errors.New("required Password")
		}
		return nil
	default:
		// Create case
		if u.Username == "" {
			return errors.New("required Username")
		}
		if u.Password == "" {
			return errors.New("required Password")
		}
		return nil
	}
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	var err = db.Debug().Create(&u).Error

	if err != nil {
		return &User{}, err
	}

	return u, nil
}

func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	users := []User{}
	var err = db.Debug().Model(&User{}).Find(&users).Error

	if err != nil {
		return &[]User{}, err
	}

	return &users, err
}

func (u *User) FindUserByUsername(db *gorm.DB, username string) (*User, error) {
	var err = db.Debug().Where(&User{Username: username}).Take(&u).Error

	if err != nil {
		return &User{}, err
	}

	return u, err
}

func (u *User) UpdateAUser(db *gorm.DB, username string) (*User, error) {

	// To hash the password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}

	db = db.Debug().Model(&User{}).Where(&User{Username: username}).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":  u.Password,
			"update_at": time.Now(),
		},
	)

	if db.Error != nil {
		return &User{}, db.Error
	}

	// This is the display the updated user
	err = db.Debug().Model(&User{}).Where(&User{Username: username}).Take(&u).Error

	if err != nil {
		return &User{}, err
	}

	return u, nil
}

func (u *User) DeleteAUser(db *gorm.DB, username string) (int64, error) {
	db = db.Debug().Model(&User{}).Where(&User{Username: username}).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

// Realm.

type Realm struct {
	gorm.Model
	Name string
}

type UserRealm struct {
	// User <> Realm
	gorm.Model
	User  User
	Realm Realm
}
