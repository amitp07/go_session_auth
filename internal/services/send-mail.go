package services

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type Email struct {
	ToAddress string
	Subject   string
	Body      string
}

func NewEmail() *Email {
	return new(Email)
}

func (e *Email) Send(email string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))

	if err != nil {
		return fmt.Errorf("could not setup aws config %s", err.Error())
	}

	sesCfg := ses.NewFromConfig(cfg)

	emailInput := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{e.ToAddress},
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data: &e.Subject,
			},
			Body: &types.Body{
				Text: &types.Content{
					Data: &e.Body,
				},
			},
		},
	}

	_, err = sesCfg.SendEmail(context.TODO(), emailInput)

	if err != nil {
		return fmt.Errorf("could not send email: %s", err.Error())
	}

	return nil
}
