package loggers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fedejuret/golang-health-checker/structures"
	"io"
	"log"
	"net/http"
	"time"
)

func Slack(service structures.Service, logger structures.ServiceLogger, text string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	slackMessage := map[string]interface{}{
		"blocks": []map[string]interface{}{
			{
				"type": "section",
				"text": map[string]string{
					"type": "mrkdwn",
					"text": fmt.Sprintf("*Service Status Update*\n\n*Service URI:* `%s`\n*Response:* `%s`\n*Timestamp:* `%s`",
						service.URI, text, timestamp),
				},
			},
			{
				"type": "divider",
			},
			{
				"type": "context",
				"elements": []map[string]string{
					{
						"type": "mrkdwn",
						"text": ":information_source: This is an automated message from the Health Checker service.",
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(slackMessage)
	if err != nil {
		panic(err.Error())
	}

	req, err := http.NewRequest("POST", logger.Webhook, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Println("Unexpected response from Slack: " + resp.Status)
	}
}
