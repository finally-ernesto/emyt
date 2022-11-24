package dbprovider

import (
	"github.com/emyt-io/emyt/dbprovider/models"
)

func (m *manager) bootstrap() {
	user := &models.User{
		Username: "root",
	}
	user.GeneratePassword()
	m.AddUser(user)
}
