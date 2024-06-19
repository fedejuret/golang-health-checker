package loggers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fedejuret/golang-health-checker/structures"
	"net/http"
	"time"
)

func Discord(service structures.Service, logger structures.ServiceLogger, text string, level string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	color := 64533

	if level == "error" {
		color = 16711680
	}

	embed := map[string]interface{}{
		"title":       "Service Status Update",
		"description": fmt.Sprintf("%s responded %s", service.URI, text),
		"color":       color,
		"fields": []map[string]string{
			{
				"name":  "Timestamp",
				"value": timestamp,
			},
			{
				"name":  "Service URI",
				"value": service.URI,
			},
			{
				"name":  "Response",
				"value": text,
			},
		},
	}

	discordMessage := map[string]interface{}{
		"embeds": []map[string]interface{}{embed},
	}

	jsonData, err := json.Marshal(discordMessage)
	if err != nil {
		panic(err.Error())
	}

	req, err := http.NewRequest("POST", logger.Webhook, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		fmt.Printf("Unexpected response from Discord: %s\n", resp.Status)
	}
}
