package dbprovider

import (
	"github.com/emyt-io/emyt/dbprovider/models"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
)

var UserManager UserManagerInterface

type UserManagerInterface interface {
	AddUser(user *models.User) error
	FindAll(users *[]models.User) error
	bootstrap() error
	// Add other methods
}

type manager struct {
	db *gorm.DB
}

func Init() {
	// Start Database Connection
	db, err := gorm.Open("sqlite3", "emyt.db")
	if err != nil {
		log.Fatal("Failed to init db:", err)
	}
	// AutoMigrate Database Models
	db.AutoMigrate(&models.User{})
	// Create Table Managers
	UserManager = &manager{db: db}
	// Bootstrap root
	err = UserManager.bootstrap()
	if err != nil {
		return
	}
}

func (mgr manager) bootstrap() (err error) {
	user := &models.User{
		Username: "root",
	}
	user.GeneratePassword()
	var users []models.User
	err = manager.FindAll(mgr, &users)
	if err != nil {
		// Handle something
	} else {
		if len(users) == 0 {
			err := mgr.AddUser(user)
			if err != nil {
				// TODO: ADD LOG
			}
		} else {
			// TODO: ADD LOG
		}
	}
	return
}

func (mgr manager) FindAll(users *[]models.User) (err error) {
	mgr.db.Find(users)
	return
}

func (mgr manager) AddUser(user *models.User) (err error) {
	mgr.db.Create(user)
	if errs := mgr.db.GetErrors(); len(errs) > 0 {
		err = errs[0]
	}
	return
}
