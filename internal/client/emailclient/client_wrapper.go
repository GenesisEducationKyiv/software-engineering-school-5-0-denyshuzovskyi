package emailclient

import (
	"context"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
	"github.com/mailgun/mailgun-go/v4"
)

type EmailClientWrapper struct {
	client *mailgun.MailgunImpl
}

func NewEmailClient(client *mailgun.MailgunImpl) *EmailClientWrapper {
	return &EmailClientWrapper{
		client: client,
	}
}

func (w *EmailClientWrapper) Send(ctx context.Context, email dto.SimpleEmail) error {
	m := mailgun.NewMessage(
		email.From,
		email.Subject,
		email.Text,
		email.To,
	)

	_, _, err := w.client.Send(ctx, m)

	return err
}
