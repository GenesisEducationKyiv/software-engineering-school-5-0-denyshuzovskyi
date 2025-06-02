package test

import (
	"context"
	"fmt"
	"github.com/mailgun/mailgun-go/v4"
	"os"
	"testing"
	"time"
)

func TestSendEmail(t *testing.T) {
	t.Skip("Demo test")
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		apiKey = ""
	}

	msg, id, err := SendSimpleMessage("", apiKey)
	fmt.Println(msg)
	fmt.Println(id)
	fmt.Println(err)
}

func SendSimpleMessage(domain, apiKey string) (string, string, error) {
	mg := mailgun.NewMailgun(domain, apiKey)
	//When you have an EU-domain, you must specify the endpoint:
	// mg.SetAPIBase("https://api.eu.mailgun.net")
	m := mailgun.NewMessage(
		"Mailgun Sandbox <postmaster@sandbox.mailgun.org>",
		"Hello",
		"Congratulations, you just sent an email with Mailgun! You are truly awesome!",
		"Denys Huzovskyi <guderu@gmail.com>",
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	msg, id, err := mg.Send(ctx, m)
	return msg, id, err
}
