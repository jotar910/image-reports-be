package endpoint

import (
	"context"
	"time"

	"image-reports/helpers/services/kafka"
)

func OnReportCreatedMessage(ctx context.Context, message *kafka.ReportCreatedMessage) (*kafka.ImageStoredMessage, error) {
	time.Sleep(time.Second * 200)
	return kafka.NewImageStoredMessageCompleted(message.ReportId, message.ImageId), nil
}

func OnReportDeletedMessage(ctx context.Context, message *kafka.ReportDeletedMessage) error {
	time.Sleep(time.Second * 200)
	return nil
}
