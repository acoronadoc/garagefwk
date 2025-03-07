package garagefwk

import (
	"os"

	"gopkg.in/yaml.v2"

	_ "github.com/go-sql-driver/mysql"
)

type Screen struct {
	Url     string                 `yaml:"url"`
	Scrtype string                 `yaml:"scrtype"`
	Options map[string]interface{} `yaml:"options"`
}

type SubMenu struct {
	Title  string     `yaml:"title"`
	Childs []MenuItem `yaml:"childs"`
}

type MenuItem struct {
	Href  string `yaml:"href"`
	Title string `yaml:"title"`
}

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Db       string `yaml:"db"`
	} `yaml:"database"`
	Menus   map[string][]SubMenu `yaml:"menus"`
	Screens []Screen             `yaml:"screens"`
}

func readConfig() Config {

	yamlFile, err := os.ReadFile("config.yaml")

	if err != nil {
		panic(err)
	}

	config := Config{}
	err = yaml.Unmarshal(yamlFile, &config)

	if err != nil {
		panic(err)
	}

	return config
}

func getScreenByUrl(config Config, url string) *Screen {

	for _, scr := range config.Screens {
		if getVarsURL(scr.Url, url) != nil {
			return &scr
		}
	}

	return nil
}
