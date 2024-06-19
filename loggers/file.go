package loggers

import (
	"fmt"
	"github.com/fedejuret/golang-health-checker/structures"
	"os"
	"path/filepath"
	"time"
)

func File(service structures.Service, logger structures.ServiceLogger, response string) {
	dir := filepath.Dir(logger.Path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err.Error())
		}
	}

	file, err := os.OpenFile(logger.Path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("%s - %s responded %s\n", timestamp, service.URI, response)

	_, err = file.WriteString(logEntry)
	if err != nil {
		panic(err.Error())
	}
}
