package mailingservices

import (
	"context"
	"github.com/mailgun/mailgun-go/v4"
	"os"
	"time"
)

func SendSimpleMessage(UserEmail, EmailSubject, EmailBody string) (string, error) {
	domain := os.Getenv("MG_DOMAIN")
	apiKey := os.Getenv("MG_PUBLIC_API_KEY")
	EmailFrom := os.Getenv("MG_EMAIL_FROM")

	mg := mailgun.NewMailgun(domain, apiKey)
	m := mg.NewMessage(EmailFrom, EmailSubject, EmailBody, UserEmail)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	res, _, err := mg.Send(ctx, m)
	if err != nil {
		return "", err
	}
	return res, nil
}
func SendVerifyAccount(userEmail, link string) (string, error) {
	domain := os.Getenv("MG_DOMAIN")
	apikey := os.Getenv("MG_PUBLIC_API_KEY")
	EmailFrom := os.Getenv("MG_EMAIL_FROM")

	mg := mailgun.NewMailgun(domain, apikey)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	m := mg.NewMessage(EmailFrom, "Verify Account", "")
	m.SetTemplate("verify.account")
	if err := m.AddRecipient(userEmail); err != nil {
		return "", err
	}

	err := m.AddVariable("link", link)
	if err != nil {
		return "", err
	}

	res, _, errr := mg.Send(ctx, m)
	return res, errr
}

func SendResetPassword(userEmail, link string) (string, error) {
	domain := os.Getenv("MG_DOMAIN")
	apikey := os.Getenv("MG_PUBLIC_API_KEY")
	EmailFrom := os.Getenv("MG_EMAIL_FROM")

	mg := mailgun.NewMailgun(domain, apikey)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	m := mg.NewMessage(EmailFrom, "Reset Password", "")
	m.SetTemplate("reset.password")
	if err := m.AddRecipient(userEmail); err != nil {
		return "", err
	}

	err := m.AddVariable("link", link)
	if err != nil {
		return "", err
	}

	res, _, errr := mg.Send(ctx, m)
	if errr != nil {
		return "", errr
	}
	return res, nil
}
