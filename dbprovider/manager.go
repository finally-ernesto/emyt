package dbprovider

import (
	"github.com/emyt-io/emyt/dbprovider/models"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
)

var Mgr Manager

type Manager interface {
	AddUser(user *models.User) error
	// Add other methods
}

type manager struct {
	db *gorm.DB
}

func init() {
	db, err := gorm.Open("sqlite3", "emyt.db")
	if err != nil {
		log.Fatal("Failed to init db:", err)
	}
	Mgr = &manager{db: db}
}

func (m *manager) AddUser(article *models.User) (err error) {
	m.db.Create(article)
	if errs := m.db.GetErrors(); len(errs) > 0 {
		err = errs[0]
	}
	return
}
