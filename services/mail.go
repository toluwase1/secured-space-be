package services

import (
	"fmt"
	m "github.com/keighl/mandrill"
	"github.com/mattbaird/gochimp"
	"log"
	"os"
)

func SendEmail(userEmail, Subject, message string) error {
	client := m.ClientWithKey(os.Getenv("MANDRILL_API"))

	mail := &m.Message{}
	mail.AddRecipient(userEmail, "", "to")
	mail.FromEmail = os.Getenv("RENTALS_EMAIL")
	mail.FromName = "Rentals"
	mail.Subject = Subject
	mail.Text = message

	_, err := client.MessagesSend(mail)
	if err != nil {
		return err
	}
	return nil
}
func Mail() {
	apiKey := os.Getenv("MANDRILL_API")
	mandrillApi, err := gochimp.NewMandrill(apiKey)
	if err != nil {
		fmt.Println("Error instantiating client")
	}
	recipients := []gochimp.Recipient{
		gochimp.Recipient{Email: "mail@securespace.ng", Type: "to"},
	}

	message := gochimp.Message{
		Subject:   "Welcome aboard!",
		FromEmail: "mail@securespace.ng",
		FromName:  "Boss Man",
		To:        recipients,
	}

	res, err := mandrillApi.MessageSend(message, false)

	if err != nil {
		fmt.Println("Error sending message")
	}
	log.Println(res)
}
