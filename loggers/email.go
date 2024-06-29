package loggers

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/fedejuret/golang-health-checker/structures"
	gomail "gopkg.in/mail.v2"
)

func Email(service structures.Service, logger structures.ServiceLogger, response string) {

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("%s - %s responded %s\n", timestamp, service.URI, response)

	m := gomail.NewMessage()

	for _, email := range logger.To {
		m.SetHeader("From", logger.SmtpConfig.Username)
		m.SetHeader("To", email)
		m.SetHeader("Subject", "Health check of "+service.URI)

		m.SetBody("text/plain", logEntry)

		d := gomail.NewDialer(logger.SmtpConfig.Host, logger.SmtpConfig.Port, logger.SmtpConfig.Username, logger.SmtpConfig.Password)

		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

		if err := d.DialAndSend(m); err != nil {
			log.Println(err)
		}
	}

}
