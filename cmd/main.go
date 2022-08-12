package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/fomik2/ticket-system/internal/app"
	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		HTTP_port string `yaml:"http_port"`
		Index     string `yaml:"templ_index"`
		Editor    string `yaml:"templ_editor"`
		Tickets   string `yaml:"tickets"`
		Counter   string `yaml:"counter"`
		CSS_path  string `yaml:"css_path"`
	}
)

func NewConfig() map[string]string {
	config := make(map[string]string)
	cfg := &Config{}
	data, err := os.Open("./config/config.yaml")
	if err != nil {
		log.Panicln("Не могу открыть файл конфигурации", err)
	}
	defer data.Close()
	byteData, err := ioutil.ReadAll(data)
	if err != nil {
		log.Panicln("Не могу прочитать файл конфигурации", err)
	}
	err = yaml.Unmarshal(byteData, &cfg)
	if err != nil {
		log.Panicln("Не могу распарсить файл конфигурации", err)
	}
	config["index"] = cfg.Index
	config["editor"] = cfg.Editor
	config["tickets"] = cfg.Tickets
	config["counter"] = cfg.Counter
	config["http_port"] = cfg.HTTP_port
	config["css_path"] = cfg.CSS_path
	return config
}

func main() {
	cfg := NewConfig()
	app.Run(cfg)
}
