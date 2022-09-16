package service

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"image-reports/helpers/services/kafka"
	log "image-reports/helpers/services/logger"
)

var options [][]string = [][]string{
	{"Dark", "Murder", "Violence"},
	{"Violence"},
	{"Dark", "Frightening", "Obscene", "Sharp", "Objects"},
	{"Frightening", "Absurd", "Obscene"},
	{"Alcohol", "Drugs", "Smoking"},
	{"Drugs", "Dark"},
	{"Sharp", "Objects"},
	{"Smoking"},
	{"Drugs"},
	{"Dark", "Absurd"},
	{"Absurd, Sports"},
	{"Beauty", "Happiness"},
	{"Friends", "Travel"},
	{"Colorful", "Friends", "Sports", "Happiness"},
	{"Animals", "Love"},
	{"Sports", "Health", "Love", "Family"},
	{"Friends", "Love", "Happiness"},
	{"Love", "Family", "Friends"},
	{"Love", "Family", "Friends", "Animals", "Health"},
	{"Colorful", "Love", "Family", "Animals", "Beauty", "Friends", "Travel", "Happiness", "Sports", "Health"},
}

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
	log.Debugf("Writing on going processing message to kafka")
	if err := pa.w.Write(context.Background(), m); err != nil {
		log.Errorf("could not write message on report created: %v", err)
	}

	time.Sleep(time.Duration(randIntn(30)) * time.Second)

	value := randIntn(100)

	if firstThreshold := randInti(50, 100); value < firstThreshold {
		index := randInti(0, len(options))
		pa.onExecuteSuccess(index*5, options[index])
		return
	}
	pa.onExecuteError(errors.New("unable to process image"))
}

func (pa *processAlgorithm) onExecuteSuccess(grade int, categories []string) {
	log.Debugf("Writing completed processing message to kafka")
	m := kafka.NewImageProcessedMessageCompleted(pa.reportId, pa.imageId, grade, categories)
	if err := pa.w.Write(context.Background(), m); err != nil {
		log.Errorf("could not write message on report created: %v", err)
	}
}

func (pa *processAlgorithm) onExecuteError(err error) {
	log.Debugf("Writing failed processing message to kafka")
	m := kafka.NewImageProcessedMessageFailed(pa.reportId, pa.imageId, err)
	if err := pa.w.Write(context.Background(), m); err != nil {
		log.Errorf("could not write message on report created: %v", err)
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
