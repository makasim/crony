package cfg

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"log"
)

type Config struct {
	Tasks []Task `json:"tasks"`
}

type Task struct {
	Name string `json:"name"`
	Cron string `json:"cron"`
	Url string `json:"url"`
}

func Load() *Config {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dat, err := ioutil.ReadFile(dir + "/crony.json")
	if err != nil {
		log.Fatal(err)
	}

	cfg := &Config{}
	if err := json.Unmarshal(dat, cfg); err != nil {
		log.Fatal(err)
	}

	return cfg
}