package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/robfig/cron"
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

func main() {
	var wg sync.WaitGroup
	config, err := Load()
	if err != nil {
		panic(err)
	}

	logOption := cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)))

	var c *cron.Cron
	if config.Location != "" {
		log.Printf("set cron location to: %s\n", config.Location)
		c = cron.New(cron.WithLocation(time.UTC), logOption)
	} else {
		c = cron.New(logOption)
	}

	for _, task := range config.Tasks {
		log.Printf("register %s %s %s\n", task.Name, task.Cron, task.Url)
		_, err := c.AddFunc(task.Cron, createCmd(wg, task))

		if err != nil {
			log.Fatalf("failed to add cron func for %s %s %s\n", task.Name, task.Cron, task.Url)
		}
	}

	wg.Add(1)
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, os.Kill)

		<-signals

		c.Stop()
		wg.Done()
	}()

	c.Start()

	log.Println("crony is working")
	wg.Wait()
}

func createCmd(wg sync.WaitGroup, task Task) func() {
	return func() {
		wg.Add(1)
		defer wg.Done()

		buf := bytes.NewBuffer([]byte{})

		resp, err := http.Post(task.Url, "text/html", buf)
		if err != nil {
			log.Printf("run %s %s %s - %s\n", task.Name, task.Cron, task.Url, err)

			return
		}

		log.Printf("run %s %s %s - response status %d\n", task.Name, task.Cron, task.Url, resp.StatusCode)
	}
}
