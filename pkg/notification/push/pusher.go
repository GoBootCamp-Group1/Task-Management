package push

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"google.golang.org/api/option"
)

type PusherPushNotifier struct {
	IPusher
	AuthKey string
}

func (p *PusherPushNotifier) Send(tokens []string, input *PushNotificationInput) error {
	// Path to your service account key file
	opt := option.WithCredentialsFile("path/to/serviceAccountKey.json")

	// Initialize the Pusher app
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return fmt.Errorf("error initializing app: %v\n", err)
	}

	// Get a messaging client from the app
	client, err := app.Messaging(context.Background())
	if err != nil {
		return fmt.Errorf("error getting Messaging client: %v\n", err)
	}

	// Create a message to send
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: "Hello, World!",
			Body:  "This is a Pusher Cloud Messaging test message!",
		},
		Token: "YOUR_DEVICE_FCM_TOKEN",
	}

	// Send the message
	response, err := client.Send(context.Background(), message)
	if err != nil {
		return fmt.Errorf("error sending message: %v\n", err)
	}

	fmt.Printf("Successfully sent message: %s\n", response)

	return nil
}
