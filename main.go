package main

import (
	"bytes"
	"github.com/formapro/crony/cfg"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"github.com/robfig/cron"
)

func main() {
	var wg sync.WaitGroup
	config := cfg.Load()
	c := cron.New()

	for _, task := range config.Tasks {
		err := c.AddFunc(task.Cron, func() {
			wg.Add(1)
			defer wg.Done()

			buf := bytes.NewBuffer([]byte{})

			resp, err := http.Post(task.Url, "text/html", buf)
			if err != nil {
				log.Printf("Run %s %s %s - %s\n", task.Name, task.Cron, task.Url, err)

				return
			}

			log.Printf("Run %s %s %s - response status %d\n", task.Name, task.Cron, task.Url, resp.StatusCode)
		})

		if err != nil {
			log.Printf("Failed to add cron func for %s %s %s\n", task.Name, task.Cron, task.Url)
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

	log.Println("Crony is working")
	wg.Wait()
}
