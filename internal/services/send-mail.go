package services

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type emailConfig struct {
	toAddress string
	subject   string
	body      string
}

func OtpEmailConfig() *emailConfig {
	subject := "Go Session Auth: OTP"

	return &emailConfig{
		subject: subject,
	}
}

// send email sychrounous process
// later add this in goroutine to run concurrently
func (e *emailConfig) Send(email string, otp string) error {

	// set body and toAddress
	e.body = fmt.Sprintf(`Hello,
Please use %s to verify your authentication.
	`, otp)

	e.toAddress = email

	accessKey := os.Getenv("SES_ACCESS_KEY")
	secretKey := os.Getenv("SES_ACCESS_SECRET")

	fmt.Printf("access %s, secret %s\n", accessKey, secretKey)
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
	)

	if err != nil {
		return fmt.Errorf("could not setup aws config %s", err.Error())
	}

	sesCfg := ses.NewFromConfig(cfg)

	fromAddr := os.Getenv("SES_FROM_ADDRESS")

	emailInput := &ses.SendEmailInput{
		Source: &fromAddr,
		Destination: &types.Destination{
			ToAddresses: []string{e.toAddress},
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data: &e.subject,
			},
			Body: &types.Body{
				Text: &types.Content{
					Data: &e.body,
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
