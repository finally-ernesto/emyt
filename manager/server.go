package manager

/*
This version using `DbOps` as separate object that handle operation to database.
This can be expanded as a core functionality to separate handler and model/logic
and encourage separation of core business logic and transport.
*/

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

func (d DbOps) findAll(users *[]User) error {
	return d.db.Find(users).Error
}

func (d DbOps) create(user *User) error {
	return d.db.Create(user).Error
}

func (d DbOps) findByPage(users *[]User, page, view int) error {
	return d.db.Limit(view).Offset(view * (page - 1)).Find(&users).Error

}

func (d DbOps) updateByName(name, email string) error {
	var user User
	d.db.Where("name=?", name).Find(&user)
	user.Email = email
	return d.db.Save(&user).Error
}

func (d DbOps) deleteByName(name string) error {
	var user User
	d.db.Where("name=?", name).Find(&user)
	return d.db.Delete(&user).Error
}

func handlerFunc(msg string) func(echo.Context) error {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, msg)
	}
}

func allUsers(dbobj DbOps) func(echo.Context) error {
	return func(c echo.Context) error {
		var users []User
		dbobj.findAll(&users)
		fmt.Println("{}", users)

		return c.JSON(http.StatusOK, users)
	}
}

func newUser(dbobj DbOps) func(echo.Context) error {
	return func(c echo.Context) error {
		name := c.Param("name")
		email := c.Param("email")
		dbobj.create(&User{Name: name, Email: email})
		return c.String(http.StatusOK, name+" user successfully created")
	}
}

func deleteUser(dbobj DbOps) func(echo.Context) error {
	return func(c echo.Context) error {
		name := c.Param("name")

		dbobj.deleteByName(name)

		return c.String(http.StatusOK, name+" user successfully deleted")
	}
}

func updateUser(dbobj DbOps) func(echo.Context) error {
	return func(c echo.Context) error {
		name := c.Param("name")
		email := c.Param("email")
		dbobj.updateByName(name, email)
		return c.String(http.StatusOK, name+" user successfully updated")
	}
}

func usersByPage(dbobj DbOps) func(echo.Context) error {
	return func(c echo.Context) error {
		limit, _ := strconv.Atoi(c.QueryParam("limit"))
		page, _ := strconv.Atoi(c.QueryParam("page"))
		var result []User
		dbobj.findByPage(&result, page, limit)
		return c.JSON(http.StatusOK, result)
	}
}

func handleRequest(dbgorm *gorm.DB) {
	e := echo.New()
	db := DbOps{dbgorm}

	e.GET("/users", allUsers(db))
	e.GET("/user", usersByPage(db))
	e.POST("/user/:name/:email", newUser(db))
	e.DELETE("/user/:name", deleteUser(db))
	e.PUT("/user/:name/:email", updateUser(db))

	e.Logger.Fatal(e.Start(":9999"))
}

func initialMigration(db *gorm.DB) {

	db.AutoMigrate(&User{})
}

func Start() {
	fmt.Println("Go ORM tutorial")
	db, err := gorm.Open("sqlite3", "emyt.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()
	initialMigration(db)
	handleRequest(db)
}
