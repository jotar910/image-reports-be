package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

var writers = make(map[string]*kafkaWriter, 0)
var readers = make(map[string]*kafkaReader, 0)

type Message = kafka.Message

type kafkaWriter struct {
	writer *kafka.Writer
}

func Writer(topic string) *kafkaWriter {
	if writer, ok := writers[topic]; ok {
		return writer
	}
	return &kafkaWriter{
		kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{brokerAddress},
			Topic:   topic,
		}),
	}
}

func (w *kafkaWriter) Write(ctx context.Context, message KafkaMessage) error {
	m, err := message.ToMessage()
	if err != nil {
		return err
	}
	return w.writer.WriteMessages(ctx, m)
}

func (w *kafkaWriter) Close() error {
	err := w.writer.Close()
	if err == nil {
		delete(writers, w.writer.Topic)
	}
	return err
}

type kafkaReader struct {
	reader *kafka.Reader
}

func Reader(topic string, group string) *kafkaReader {
	if reader, ok := readers[topic]; ok {
		return reader
	}
	return &kafkaReader{
		kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{brokerAddress},
			Topic:   topic,
			GroupID: group,
		}),
	}
}

func (r *kafkaReader) Read(ctx context.Context, message KafkaMessage) error {
	m, err := r.reader.ReadMessage(ctx)
	if err != nil {
		return err
	}
	return message.FromMessage(m)
}

func (r *kafkaReader) Close() error {
	err := r.reader.Close()
	if err == nil {
		delete(readers, r.reader.Config().Topic)
	}
	return err
}

func Shutdown() {
	for _, writer := range writers {
		if err := writer.Close(); err != nil {
			log.Println(fmt.Errorf("closing topic %s: %w", writer.writer.Topic, err))
		}
	}
}
