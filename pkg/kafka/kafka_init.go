package kafkainit

import (
	"fmt"
	"github.com/segmentio/kafka-go"
)

/*
func checkContextKafka(goodsproducer GoodsProducer, ctx context.Context) error {
	<-ctx.Done()

	err := goodsproducer.writer.Close()

	if err != nil {
		fmt.Errorf("не удалось закрыть writer: %w", err)
	}

	return nil
}
*/

var (
	kafkaBroker      = "kafka:9092"
	kafkaTopic       = "new_goods_topic."
	groupID          = "consumer-group-1"
	kafkaPartition   = 1
	kafkaReplication = 1
)

func createTopic() error {

	conn, err := kafka.Dial("tcp", kafkaBroker)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к Kafka: %w", err)
	}
	defer conn.Close()

	topics, err := conn.ReadPartitions(kafkaTopic)
	if err == nil && len(topics) > 0 {
		return nil
	}

	topicConfig := kafka.TopicConfig{
		Topic:             kafkaTopic,
		NumPartitions:     kafkaPartition,
		ReplicationFactor: kafkaReplication,
	}

	err = conn.CreateTopics(topicConfig)
	if err != nil {
		return fmt.Errorf("не удалось создать топик: %w", err)
	}

	return nil
}

func NewKafkaReader() *kafka.Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaBroker},
		Topic:    kafkaTopic,
		GroupID:  groupID,
		MinBytes: 1e3,  // 1KB
		MaxBytes: 10e6, // 10MB
	})

	return reader
}

func NewKafkaWriter() (*kafka.Writer, error) {

	err := createTopic()

	if err != nil {
		return nil, fmt.Errorf("Ошибка при создании топика: %w", err)
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    kafkaTopic,
		Balancer: &kafka.LeastBytes{},
	}

	return writer, nil
}
