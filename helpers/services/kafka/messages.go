package kafka

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"strconv"

	"github.com/segmentio/kafka-go"
)

type KafkaMessage interface {
	ToMessage() (Message, error)
	FromMessage(source Message) error
}

type kafkaMessage struct {
	ReportId int
	ImageId  string
}

func (m *kafkaMessage) ToMessage() (Message, error) {
	value, err := getBytes(m)
	if err != nil {
		return kafka.Message{}, err
	}
	return Message{
		Key:   []byte(strconv.Itoa(m.ReportId)),
		Value: value,
	}, nil
}

func (m *kafkaMessage) FromMessage(source Message) error {
	decoder := json.NewDecoder(bytes.NewReader(source.Value))
	return decoder.Decode(m)
}

type reportCreatedMessage struct {
	kafkaMessage
	ReportImage []byte
}

func NewReportCreatedMessage(reportId int, imageId string, reportImage []byte) *reportCreatedMessage {
	return &reportCreatedMessage{
		kafkaMessage: kafkaMessage{reportId, imageId},
		ReportImage:  reportImage,
	}
}

type deleteCreatedMessage struct {
	kafkaMessage
}

func NewDeletedReportMessage(reportId int, imageId string) *deleteCreatedMessage {
	return &deleteCreatedMessage{
		kafkaMessage: kafkaMessage{reportId, imageId},
	}
}

type imageProcessedMessage struct {
	kafkaMessage
	Grade      int
	Categories []string
	Err        error
}

func NewImageProcessedMessageCompleted(reportId int, imageId string, grade int, categories []string) *imageProcessedMessage {
	return &imageProcessedMessage{
		kafkaMessage: kafkaMessage{reportId, imageId},
		Grade:        grade,
		Categories:   categories,
	}
}

func NewImageProcessedMessageFailed(reportId int, imageId string, err error) *imageProcessedMessage {
	return &imageProcessedMessage{
		kafkaMessage: kafkaMessage{reportId, imageId},
		Err:          err,
	}
}

type imageStoredMessage struct {
	kafkaMessage
	Err error
}

func NewImageStoredMessageCompleted(reportId int, imageId string) *imageStoredMessage {
	return &imageStoredMessage{
		kafkaMessage: kafkaMessage{reportId, imageId},
	}
}

func NewImageStoredMessageFailed(reportId int, imageId string, err error) *imageStoredMessage {
	return &imageStoredMessage{
		kafkaMessage: kafkaMessage{reportId, imageId},
		Err:          err,
	}
}

func getBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
