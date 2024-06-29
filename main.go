package main

import (
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"

	"github.com/fedejuret/golang-health-checker/loggers"
	"github.com/fedejuret/golang-health-checker/structures"
	"github.com/jasonlvhit/gocron"
)

func main() {
	var wg sync.WaitGroup

	err := filepath.Walk("services", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}

		if info.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}

		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			processServiceFile(path)
		}(path)

		return nil
	})

	if err != nil {
		log.Println(err)
	}

	wg.Wait()
	<-gocron.Start()
}

func processServiceFile(path string) {
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Println(err)
		return
	}

	var requestStructure structures.Service
	err = json.Unmarshal(byteValue, &requestStructure)
	if err != nil {
		log.Println(err)
		return
	}

	registerCronJobs(requestStructure)
}

func registerCronJobs(service structures.Service) {
	gocron.Every(uint64(service.Every)).Second().Do(func() {
		go checkService(service)
	})
}

func checkService(service structures.Service) {
	log.Println("Checking " + service.URI + " ...")

	request, err := http.NewRequest("GET", service.URI, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}

	for _, header := range service.Headers {
		request.Header.Set(header.Key, header.Value)
	}

	client := &http.Client{
		Timeout: time.Duration(service.Timeout) * time.Second,
	}
	response, err := client.Do(request)
	if err != nil {
		log.Println("Error making request:", err)
		for _, logger := range service.Loggers {
			go dispatchNotification(service, logger, "Error making request: "+err.Error(), "error")
		}
		return
	}
	defer response.Body.Close()

	success := slices.Contains(service.AcceptedHTTPStatusCodes, response.StatusCode)
	level := "error"
	if success {
		level = "success"
	}

	for _, logger := range service.Loggers {
		if slices.Contains(logger.Level, level) {
			go dispatchNotification(service, logger, response.Status, level)
		}
	}
}

func dispatchNotification(service structures.Service, logger structures.ServiceLogger, text string, level string) {
	switch logger.Type {
	case "file":
		go loggers.File(service, logger, text)
	case "discord":
		go loggers.Discord(service, logger, text, level)
	case "slack":
		go loggers.Slack(service, logger, text)
	case "email":
		go loggers.Email(service, logger, text)
	}
}
