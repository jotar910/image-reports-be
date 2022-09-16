package service

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"image-reports/helpers/services/kafka"
	log "image-reports/helpers/services/logger"
)

type processAlgorithm struct {
	reportId uint
	imageId  string
	w        *kafka.KafkaWriter
}

func newProcessAlgorithm(reportId uint, imageId string) *processAlgorithm {
	return &processAlgorithm{reportId, imageId, kafka.Writer(kafka.TopicImageProcessed)}
}

func (pa *processAlgorithm) execute() {
	m := kafka.NewImageProcessedMessageGoing(pa.reportId)
	if err := pa.w.Write(context.Background(), m); err != nil {
		log.Errorf("could not write message on report created: %w", err)
	}

	time.Sleep(time.Duration(randIntn(30)) * time.Second)

	value := randIntn(100)

	if firstThreshold := randInti(50, 100); value < firstThreshold {
		pa.onExecuteSuccess(0, []string{}) // TODO
		return
	}
	pa.onExecuteError(errors.New("unable to process image"))
}

func (pa *processAlgorithm) onExecuteSuccess(grade int, categories []string) {
	w := kafka.Writer(kafka.TopicImageProcessed)
	m := kafka.NewImageProcessedMessageCompleted(pa.reportId, pa.imageId, grade, categories)
	if err := w.Write(context.Background(), m); err != nil {
		log.Errorf("could not write message on report created: %w", err)
	}
}

func (pa *processAlgorithm) onExecuteError(err error) {
	w := kafka.Writer(kafka.TopicImageProcessed)
	m := kafka.NewImageProcessedMessageFailed(pa.reportId, pa.imageId, err)
	if err := w.Write(context.Background(), m); err != nil {
		log.Errorf("could not write message on report created: %w", err)
	}
}

func randIntn(max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max)
}

func randInti(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}
