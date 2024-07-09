package notification

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"gopkg.in/gomail.v2"
	_ "gopkg.in/gomail.v2"
	"html/template"
	"os"
	"sync"
)

type EmailNotifierConf struct {
	NotifierConf
	// Add necessary fields, like SMTP server info
	SmtpHost        string
	SmtpPort        int
	SmtpUsername    string
	SmtpPassword    string
	SmtpFromAddress string
	SmtpEncryption  string
	SmtpFromName    string
}

func NewEmailNotifier(cfg *EmailNotifierConf) (*EmailNotifierConf, error) {
	return &EmailNotifierConf{
		SmtpHost:        cfg.SmtpHost,
		SmtpPort:        cfg.SmtpPort,
		SmtpUsername:    cfg.SmtpUsername,
		SmtpPassword:    cfg.SmtpPassword,
		SmtpFromAddress: cfg.SmtpFromAddress,
		SmtpEncryption:  cfg.SmtpEncryption,
		SmtpFromName:    cfg.SmtpFromName,
	}, nil
}

func (e *EmailNotifierConf) RenderTemplateToString(templatePath string, data any) (string, error) {
	// Read the template file
	tmplContent, err := os.ReadFile(templatePath)
	if err != nil {
		return "", err
	}

	// Parse the template
	tmpl, err := template.New("email").Parse(string(tmplContent))
	if err != nil {
		return "", err
	}

	// Create a buffer to hold the output of the template execution
	var buf bytes.Buffer

	// Execute the template and inject the data
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	// Convert the buffer to a string and return it
	return buf.String(), nil
}

func (e *EmailNotifierConf) Send(to string, subject string, message string) error {
	// Implement email sending logic
	fmt.Printf("Sending email to %s\n", to)

	m := gomail.NewMessage()
	m.SetHeader("From", e.SmtpFromAddress)
	m.SetHeader("To", to)

	buff := make([]byte, 8)
	_, err := rand.Read(buff)
	if err != nil {
		return err
	}
	randomStr := base64.StdEncoding.EncodeToString(buff)

	m.SetHeader("Message-Id", "<"+randomStr+"@theQame.com>")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", message+"</br>")

	d := gomail.NewDialer(e.SmtpHost, e.SmtpPort, e.SmtpUsername, e.SmtpPassword)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	fmt.Println("Email Sent Successfully!")

	return nil
}

func (e *EmailNotifierConf) SendBatch(recipients []string, subject string, message string) error {
	//TODO: using worker pool
	var wg sync.WaitGroup

	for _, emailAddress := range recipients {
		wg.Add(1)

		go func() {
			defer wg.Done()
			err := e.Send(emailAddress, subject, message)
			if err != nil {
				//TODO: handle in database or other place!!!

			}
		}()
	}

	wg.Wait()

	return nil
}

func (e *EmailNotifierConf) SendCustom() {
	panic("dssadsadas")
}
