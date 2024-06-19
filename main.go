package main

import (
	"encoding/json"
	"github.com/fedejuret/golang-health-checker/structures"
	"github.com/jasonlvhit/gocron"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	err := filepath.Walk("services", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		jsonFile, err := os.Open(path)
		if err != nil {
			log.Println(err)
			return err
		}
		defer func(jsonFile *os.File) {
			err := jsonFile.Close()
			if err != nil {
				log.Println(err)
			}
		}(jsonFile)

		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			log.Println(err)
			return err
		}

		var requestStructure structures.Service
		err = json.Unmarshal(byteValue, &requestStructure)
		if err != nil {
			log.Println(err)
			return err
		}

		registerCronJobs(requestStructure)
		return nil
	})

	if err != nil {
		log.Println(err)
	}

	<-gocron.Start()
}

func registerCronJobs(service structures.Service) {
	gocron.Every(uint64(service.Every)).Second().Do(func() {
		checkService(service)
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
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Error closing response body:", err)
		}
	}(response.Body)

	if !containsStatusCode(service.AcceptedHTTPStatusCodes, response.StatusCode) {
		log.Println("Request returned wrong status code")
	}

	log.Println("Service: " + service.URI + " returned status code: " + strconv.Itoa(response.StatusCode))
}

func containsStatusCode(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
