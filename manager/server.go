package manager

/*
This version using `DbOps` as separate object that handle operation to database.
This can be expanded as a core functionality to separate handler and model/logic
and encourage separation of core business logic and transport.
*/

import (
	"fmt"
	dbprovider "github.com/emyt-io/emyt/dbprovider"
	dbModels "github.com/emyt-io/emyt/dbprovider/models"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	"net/http"
)

func handlerFunc(msg string) func(echo.Context) error {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, msg)
	}
}

func allUsers(manager dbprovider.UserManagerInterface) func(echo.Context) error {
	return func(c echo.Context) error {
		var users []dbModels.User
		err := manager.FindAll(&users)
		if err != nil {
			return err
		}
		fmt.Println("{}", users)
		return c.JSON(http.StatusOK, users)
	}
}

func handleRequest() {
	e := echo.New()
	manager := dbprovider.UserManager
	e.GET("/users", allUsers(manager))
	e.Logger.Fatal(e.Start(":9999"))
}

func Start() {
	handleRequest()
}
