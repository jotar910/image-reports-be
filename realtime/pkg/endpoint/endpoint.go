package endpoint

import (
	"context"
	"time"

	"image-reports/helpers/services/kafka"
)

func OnImageProcessedMessage(ctx context.Context, message *kafka.ImageProcessedMessage) error {
	time.Sleep(time.Second * 200)
	return nil
}

func OnImageStoredMessage(ctx context.Context, message *kafka.ImageStoredMessage) error {
	time.Sleep(time.Second * 200)
	return nil
}
