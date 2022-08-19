package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/fomik2/ticket-system/internal/app"
	rep "github.com/fomik2/ticket-system/internal/repo"
	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		HTTP_port string `yaml:"http_port"`
		Index     string `yaml:"templ_index"`
		Layout    string `yaml:"templ_layout"`
		Editor    string `yaml:"templ_editor"`
		Tickets   string `yaml:"tickets"`
		Counter   string `yaml:"counter"`
		CSS_path  string `yaml:"css_path"`
	}
)

func NewConfig() (index, layout, editor, tickets, counter, http_port, css_path string) {

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
	tickets = cfg.Tickets
	counter = cfg.Counter
	http_port = cfg.HTTP_port
	css_path = cfg.CSS_path
	return
}

func main() {
	index, layout, editor, tickets, counter, http_port, css_path := NewConfig()
	repo := rep.New(tickets, counter)
	err := app.Run(index, layout, editor, tickets, counter, http_port, css_path, repo)
	if err != nil {
		log.Println("Problem related to starting server", err)
		return
	}
}
