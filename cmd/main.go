package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/emyt-io/emyt/config"
	"github.com/emyt-io/emyt/db"
	"github.com/emyt-io/emyt/models"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/natefinch/lumberjack.v2"
)

const AppYamlFilename = "app.yaml"

func load() config.Config {
	var cfg config.Config
	// read configuration from the file and environment variables
	if err := cleanenv.ReadConfig(AppYamlFilename, &cfg); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	return cfg
}

var hosts = map[string]*models.Host{}
var redirectUrls = map[string]*models.RedirectUrl{}

func main() {

	// Load ENV
	cfg := load()
	// Init DB
	db.Init()
	// Hosts
	for _, service := range cfg.Services {
		// Service Target
		tenant := echo.New()
		var targets []*middleware.ProxyTarget
		skip := !service.UseAuth

		tenant.Use(middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
			Skipper: func(ctx echo.Context) bool { return skip },
			Validator: func(username, password string, ctx echo.Context) (bool, error) {
				return handleLogIn(username, password)
			},
		}))

		// Service Config
		if service.Type == "proxy" {
			// Web endpoint
			urlS, err := url.Parse(service.EgressUrl)
			if err != nil {
				tenant.Logger.Fatal(err)
			}
			targets = append(targets, &middleware.ProxyTarget{
				URL: urlS,
			})
			tenant.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))
			tenant.GET("/*", func(c echo.Context) error {
				return c.String(http.StatusOK, "Tenant:"+c.Request().Host)
			})
			hosts[service.IngressUrl] = &models.Host{Echo: tenant}
		} else if service.Type == "static" {
			// Static endpoint
			tenant.Use(middleware.GzipWithConfig(middleware.GzipConfig{
				Level: 5,
			}))
			tenant.Use(expiresServerHeader)
			tenant.Use(middleware.BodyLimit("10M"))
			tenant.Use(middleware.StaticWithConfig(middleware.StaticConfig{
				Root:   service.EgressUrl,
				Browse: true,
				HTML5:  true,
			}))
			// Add to Hosts
			hosts[service.IngressUrl] = &models.Host{Echo: tenant}
		} else if service.Type == "redirect" {
			// Redirect endpoint
			urlS, err := url.Parse(service.EgressUrl)
			if err != nil {
				tenant.Logger.Fatal(err)
			}
			redirectUrls[service.IngressUrl] = &models.RedirectUrl{
				URL: urlS,
			}
			tenant.Any("/*", func(c echo.Context) error {
				redirectUrl := redirectUrls[c.Request().Host].URL.String()
				return c.Redirect(http.StatusMovedPermanently, redirectUrl)
			})
			hosts[service.IngressUrl] = &models.Host{Echo: tenant}
		}
	}

	//---------
	// ROOT
	//---------
	server := echo.New()
	server.Use(middleware.Recover())
	server.GET("/status", func(c echo.Context) error {
		return c.String(http.StatusOK, "{\"success\":\"ok\"}")
	})
	hosts[cfg.StatusHost+":"+cfg.ProxyListenPort] = &models.Host{Echo: server}

	// Server
	e := echo.New()
	e.Use(middleware.Logger())
	e.Logger.SetOutput(&lumberjack.Logger{
		Filename:   cfg.Logfile,
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	})
	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		host := hosts[req.Host]
		if host == nil {
			err = echo.ErrNotFound
			e.Logger.Info("Resource Not found - " + req.Host)
		} else {
			host.Echo.ServeHTTP(res, req)
		}

		return
	})
	// 4 Terabyte limit
	e.Use(middleware.BodyLimit("4T"))

	// Start server with Graceful Shutdown WITH CERT
	//go func() {
	//	if err := e.StartTLS(":"+cfg.SSLPort,
	//		"/etc/emty/ssl/server.crt",
	//		"/etc/emyt/ssl/server.key"); err != nil && err != http.ErrServerClosed {
	//		e.Logger.Fatal("shutting down the server")
	//	}
	//}()

	// Start server with Graceful Shutdown WITHOUT CERT
	go func() {
		if err := e.Start(":9000"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

// ServerHeader middleware adds a `Server` header to the response.
func expiresServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "public, max-age=3600")
		return next(c)
	}
}

func handleLogIn(username, password string) (bool, error) {

	user := models.User{
		Username: username,
		Password: password,
	}
	user.Prepare()
	err := user.Validate("login")

	if err != nil {
		return false, err
	}

	return signIn(user.Username, user.Password)
}

func signIn(username, password string) (bool, error) {
	var err error

	u := models.User{}
	db := db.DbManager()
	err = db.Model(models.User{}).Where(&models.User{Username: username}).Take(&u).Error

	if err != nil {
		return false, err
	}

	err = models.VerifyPassword(u.Password, password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return false, err
	}

	return true, nil
}
