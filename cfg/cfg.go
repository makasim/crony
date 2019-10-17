package cfg

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Location string `json:"location"`
	Tasks []Task `json:"tasks"`
}

type Task struct {
	Name string `json:"name"`
	Cron string `json:"cron"`
	Url string `json:"url"`
}

func Load() (cfg Config, err error) {
	dir, err := os.Getwd()
	if err != nil {
		return cfg, err
	}

	dat, err := ioutil.ReadFile(dir + "/crony.json")
	if err != nil {
		return cfg, err
	}

	if err := json.Unmarshal(dat, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
