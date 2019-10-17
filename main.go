package main

import (
	"bytes"
	"github.com/makasim/crony/cfg"
	"github.com/robfig/cron"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	config, err := cfg.Load()
	if err != nil {
		panic(err)
	}

	var c *cron.Cron
	if config.Location != "" {
		c = cron.New(cron.WithLocation(time.UTC))
	} else {
		c = cron.New()
	}

	for _, task := range config.Tasks {
		_, err := c.AddFunc(task.Cron, func() {
			wg.Add(1)
			defer wg.Done()

			buf := bytes.NewBuffer([]byte{})

			resp, err := http.Post(task.Url, "text/html", buf)
			if err != nil {
				log.Printf("run %s %s %s - %s\n", task.Name, task.Cron, task.Url, err)

				return
			}

			log.Printf("run %s %s %s - response status %d\n", task.Name, task.Cron, task.Url, resp.StatusCode)
		})

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
