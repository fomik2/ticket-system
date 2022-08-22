package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/fomik2/ticket-system/internal/app"
	rep "github.com/fomik2/ticket-system/internal/repo"
	"github.com/fomik2/ticket-system/pkg/sqlite"
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
	}
)

func NewConfig() (index, layout, editor, auth, user_create, http_port, css_path, database string) {

	cfg := &Config{}
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
	index = cfg.Index
	editor = cfg.Editor
	layout = cfg.Layout
	http_port = cfg.HTTP_port
	css_path = cfg.CSS_path
	database = cfg.Database
	user_create = cfg.UserCreate
	auth = cfg.Auth
	return
}

func main() {
	index, layout, editor, auth, user_create, http_port, css_path, database := NewConfig()
	db, err := sqlite.New(database)
	if err != nil {
		log.Println("can't connect to database", err)
		return
	}
	repo := rep.New(db)
	err = app.Run(index, layout, editor, auth, user_create, http_port, css_path, database, repo)
	if err != nil {
		log.Println("Problem related to starting server", err)
		return
	}
}
