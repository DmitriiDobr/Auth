package notification

import (
	"auth/internal/auth/types"
	"context"
	kafkaNotification "github.com/DmitriiDobr/kafkaNotification/pkg"
)

type Notification struct {
	kafkaClient *kafkaNotification.Client
}

func NewNotification(kafkaClient *kafkaNotification.Client) *Notification {
	return &Notification{kafkaClient: kafkaClient}
}

func (n *Notification) Notify(ctx context.Context, message types.Message) error {
	msg := kafkaNotification.Message{
		UserID: message.UserID,
		Status: kafkaNotification.Status(message.Status),
		Header: message.Header,
		Body:   message.Body,
	}
	return n.kafkaClient.Notify(ctx, msg)
}
