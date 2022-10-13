package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/fomik2/ticket-system/internal/handlers"
	rep "github.com/fomik2/ticket-system/internal/repo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"

	"github.com/labstack/echo-contrib/prometheus"
	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		HTTP_port  string `yaml:"http_port"`
		Index      string `yaml:"templ_index"`
		Layout     string `yaml:"templ_layout"`
		Editor     string `yaml:"templ_editor"`
		UserCreate string `yaml:"templ_user_create"`
		Auth       string `yaml:"templ_auth"`
		Counter    string `yaml:"counter"`
		CSS_path   string `yaml:"css_path"`
		Database   string `yaml:"db_file"`
		SecretKey  string `yaml:"session_and_jwt_secret"`
	}
)

func NewConfig() (Config, error) {

	cfg := Config{}
	data, err := os.Open("./config/config.yaml")
	if err != nil {
		log.Println("Не могу открыть файл конфигурации", err)
	}
	defer data.Close()
	byteData, err := ioutil.ReadAll(data)
	if err != nil {
		log.Println("Не могу прочитать файл конфигурации", err)
	}
	err = yaml.Unmarshal(byteData, &cfg)
	if err != nil {
		log.Println("Не могу распарсить файл конфигурации", err)
	}
	return cfg, err
}

// NewDBConnection connet to DB
func NewDBConnection(database string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		return nil, fmt.Errorf("can't connect to db %w", err)
	}
	return db, nil
}

func main() {
	cfg, err := NewConfig()

	index := cfg.Index
	editor := cfg.Editor
	layout := cfg.Layout
	http_port := cfg.HTTP_port
	database := cfg.Database
	user_create := cfg.UserCreate
	auth := cfg.Auth
	secret := cfg.SecretKey

	db, err := NewDBConnection(database)
	if err != nil {
		log.Println("can't connect to database", err)
		return
	}
	repo := rep.New(db)

	handler, err := handlers.New(index, layout, editor, auth, user_create, secret, repo)
	if err != nil {
		log.Println(err)
		return
	}

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	//Accesible without auth
	e.POST("/api/signin", handler.APISignin)
	e.GET("/login", handler.Login)
	e.POST("/login", handler.LoginHandler)
	e.GET("/logout", handler.LogoutHandler)

	//Prometheus metrics
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	// API restricted group with JWT

	apiGR := e.Group("/api/tickets/")
	config := middleware.JWTConfig{
		Claims:     &handlers.JWTCustomClaims{},
		SigningKey: []byte(secret),
	}
	apiGR.Use(middleware.JWTWithConfig(config))
	apiGR.GET(":id", handler.APIGetTicket)
	apiGR.GET("byuser", handler.APIGetListTicketsByUser)
	apiGR.GET("", handler.APIGetListTickets)
	apiGR.POST(":id", handler.APIUpdateTicket)
	apiGR.POST("", handler.APICreateTicket)
	apiGR.DELETE(":id", handler.APIDeleteTicket)

	//HTTP restricted group auth without JWT
	httpGR := e.Group("/")

	httpGR.Use(handler.Authentication)
	httpGR.GET("", handler.WelcomeHandler)

	httpGR.POST("", handler.CreateTicket)
	httpGR.POST("tickets/:id", handler.EditHandler)
	httpGR.POST("tickets/:id/delete/", handler.DeleteHandler)
	httpGR.GET("tickets/:id", handler.GetTicketForEdit)
	httpGR.GET("user_create/", handler.CreateUserGet)
	httpGR.POST("user_create/", handler.CreateUser)

	e.Static("/css/", "../css")
	e.Use(middleware.Static("./css"))
	e.Logger.Fatal(e.Start(http_port))
}
