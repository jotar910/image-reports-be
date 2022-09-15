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
}

func newProcessAlgorithm(reportId uint, imageId string) *processAlgorithm {
	return &processAlgorithm{reportId, imageId}
}

func (pa *processAlgorithm) execute() {
	time.Sleep(time.Duration(randIntn(30000)) * time.Second)

	value := randIntn(100)
	pa.onExecuteSuccess(0, []string{}) // TODO

	if firstThreshold := randInti(50, 100); value < firstThreshold {
		pa.onExecuteSuccess(0, []string{}) // TODO
		return
	}
	pa.onExecuteError(errors.New("unable to process image"))
}

func (pa *processAlgorithm) onExecuteSuccess(grade int, categories []string) {
	println("writing")
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
