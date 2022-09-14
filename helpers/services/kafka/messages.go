package kafka

import (
	"bytes"
	"encoding/gob"
	"strconv"

	"github.com/segmentio/kafka-go"
)

type KafkaMessage interface {
	ToMessage() (Message, error)
	FromMessage(source Message) error
}

type ReportCreatedMessage struct {
	ReportId    int
	ImageId     string
	ReportImage []byte
}

func (m *ReportCreatedMessage) ToMessage() (Message, error) {
	return toMessage(m.ReportId, m)
}

func (m *ReportCreatedMessage) FromMessage(source Message) error {
	return fromMessage(source, m)
}

func NewEmptyReportCreatedMessage() *ReportCreatedMessage {
	return &ReportCreatedMessage{}
}

func NewReportCreatedMessage(reportId int, imageId string, reportImage []byte) *ReportCreatedMessage {
	return &ReportCreatedMessage{
		ReportId:    reportId,
		ImageId:     imageId,
		ReportImage: reportImage,
	}
}

type ReportDeletedMessage struct {
	ReportId int
	ImageId  string
}

func (m *ReportDeletedMessage) ToMessage() (Message, error) {
	return toMessage(m.ReportId, m)
}

func (m *ReportDeletedMessage) FromMessage(source Message) error {
	return fromMessage(source, m)
}

func NewEmptyDeletedReportMessage() *ReportDeletedMessage {
	return &ReportDeletedMessage{}
}

func NewDeletedReportMessage(reportId int, imageId string) *ReportDeletedMessage {
	return &ReportDeletedMessage{
		ReportId: reportId,
		ImageId:  imageId,
	}
}

type ImageProcessedMessage struct {
	ReportId   int
	ImageId    string
	Grade      int
	Categories []string
	Err        error
}

func (m *ImageProcessedMessage) ToMessage() (Message, error) {
	return toMessage(m.ReportId, m)
}

func (m *ImageProcessedMessage) FromMessage(source Message) error {
	return fromMessage(source, m)
}

func NewEmptyImageProcessedMessage() *ImageProcessedMessage {
	return &ImageProcessedMessage{}
}

func NewImageProcessedMessageCompleted(reportId int, imageId string, grade int, categories []string) *ImageProcessedMessage {
	return &ImageProcessedMessage{
		ReportId:   reportId,
		ImageId:    imageId,
		Grade:      grade,
		Categories: categories,
	}
}

func NewImageProcessedMessageFailed(reportId int, imageId string, err error) *ImageProcessedMessage {
	return &ImageProcessedMessage{
		ReportId: reportId,
		ImageId:  imageId,
		Err:      err,
	}
}

type ImageStoredMessage struct {
	ReportId int
	ImageId  string
	Err      error
}

func (m *ImageStoredMessage) ToMessage() (Message, error) {
	return toMessage(m.ReportId, m)
}

func (m *ImageStoredMessage) FromMessage(source Message) error {
	return fromMessage(source, m)
}

func NewEmptyImageStoredMessage() *ImageStoredMessage {
	return &ImageStoredMessage{}
}

func NewImageStoredMessageCompleted(reportId int, imageId string) *ImageStoredMessage {
	return &ImageStoredMessage{
		ReportId: reportId,
		ImageId:  imageId,
	}
}

func NewImageStoredMessageFailed(reportId int, imageId string, err error) *ImageStoredMessage {
	return &ImageStoredMessage{
		ReportId: reportId,
		ImageId:  imageId,
		Err:      err,
	}
}

func toMessage(key int, m KafkaMessage) (Message, error) {
	value, err := getBytes(m)
	if err != nil {
		return kafka.Message{}, err
	}
	return Message{
		Key:   []byte(strconv.Itoa(key)),
		Value: value,
	}, nil
}

func fromMessage(source Message, m KafkaMessage) error {
	decoder := gob.NewDecoder(bytes.NewReader(source.Value))
	return decoder.Decode(m)
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
