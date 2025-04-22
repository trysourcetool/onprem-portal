package mail

import (
	"context"

	"github.com/resend/resend-go/v2"

	"github.com/trysourcetool/onprem-portal/internal/config"
	"github.com/trysourcetool/onprem-portal/internal/errdefs"
)

type Input struct {
	To       []string
	From     string
	FromName string
	Subject  string
	Body     string
}

func Send(ctx context.Context, input Input) error {
	client := resend.NewClient(config.Config.Resend.APIKey)
	email := resend.SendEmailRequest{
		From:    input.From,
		To:      input.To,
		Subject: input.Subject,
		Text:    input.Body,
	}

	_, err := client.Emails.SendWithContext(ctx, &email)
	if err != nil {
		return errdefs.ErrResend(err)
	}

	return nil
}
