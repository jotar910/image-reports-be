package endpoint

import (
	"context"
	"time"

	"image-reports/helpers/services/kafka"
)

func OnReportCreatedMessage(ctx context.Context, message *kafka.ReportCreatedMessage) (*kafka.ImageProcessedMessage, error) {
	time.Sleep(time.Second * 200)
	return kafka.NewImageProcessedMessageCompleted(message.ReportId, message.ImageId, 100, []string{"unicorns"}), nil
}

func OnReportDeletedMessage(ctx context.Context, message *kafka.ReportDeletedMessage) error {
	time.Sleep(time.Second * 200)
	return nil
}
