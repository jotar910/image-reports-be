package kafka

import (
	"context"

	log "image-reports/helpers/services/logger"

	"github.com/segmentio/kafka-go"
)

var writers = make(map[string]*KafkaWriter, 0)
var readers = make(map[string]*KafkaReader, 0)

type Message = kafka.Message

type KafkaWriter struct {
	writer *kafka.Writer
}

func Writer(topic string) *KafkaWriter {
	if writer, ok := writers[topic]; ok {
		return writer
	}
	return &KafkaWriter{
		kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{brokerAddress},
			Topic:   topic,
		}),
	}
}

func (w *KafkaWriter) Write(ctx context.Context, message KafkaMessage) error {
	m, err := message.ToMessage()
	if err != nil {
		return err
	}
	return w.writer.WriteMessages(ctx, m)
}

func (w *KafkaWriter) Close() error {
	err := w.writer.Close()
	if err == nil {
		delete(writers, w.writer.Topic)
	}
	return err
}

type KafkaReader struct {
	reader *kafka.Reader
}

func Reader(topic string, group string) *KafkaReader {
	if reader, ok := readers[topic]; ok {
		return reader
	}
	return &KafkaReader{
		kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{brokerAddress},
			Topic:   topic,
			GroupID: group,
		}),
	}
}

func (r *KafkaReader) Read(ctx context.Context, message KafkaMessage) error {
	m, err := r.reader.ReadMessage(ctx)
	if err != nil {
		return err
	}
	return message.FromMessage(m)
}

func (r *KafkaReader) Close() error {
	err := r.reader.Close()
	if err == nil {
		delete(readers, r.reader.Config().Topic)
	}
	return err
}

func Shutdown() {
	for _, writer := range writers {
		if err := writer.Close(); err != nil {
			log.Errorf("closing topic %s: %w", writer.writer.Topic, err)
		}
	}
}
