package firebase

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type FCMService struct {
	client *messaging.Client
	ctx    context.Context
}

func NewFCMService(credentialsPath string) (*FCMService, error) {
	ctx := context.Background()

	var app *firebase.App
	var err error

	if credentialsPath != "" {
		opt := option.WithCredentialsFile(credentialsPath)
		app, err = firebase.NewApp(ctx, nil, opt)
	} else {
		// Use default credentials (from environment)
		app, err = firebase.NewApp(ctx, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase app: %w", err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get messaging client: %w", err)
	}

	log.Println("✅ Firebase/FCM initialized")

	return &FCMService{
		client: client,
		ctx:    ctx,
	}, nil
}

// SendPushNotification sends a push notification to a specific device
func (f *FCMService) SendPushNotification(token, title, body string, data map[string]string) error {
	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
		Android: &messaging.AndroidConfig{
			Priority: "high",
		},
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority": "10",
			},
		},
	}

	response, err := f.client.Send(f.ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send FCM message: %w", err)
	}

	log.Printf("✅ Push notification sent: %s", response)
	return nil
}

// SendMulticast sends a notification to multiple devices
func (f *FCMService) SendMulticast(tokens []string, title, body string, data map[string]string) error {
	message := &messaging.MulticastMessage{
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
	}

	response, err := f.client.SendMulticast(f.ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send multicast message: %w", err)
	}

	log.Printf("✅ Multicast sent: %d success, %d failure", response.SuccessCount, response.FailureCount)
	return nil
}
