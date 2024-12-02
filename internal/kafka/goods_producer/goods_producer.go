package goods_producer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log/slog"
	"time"
)

type kafkaMessageType struct {
	goodsIds []int64
	dt       time.Time
}

type GoodsProducer struct {
	log    *slog.Logger
	writer *kafka.Writer
}

func New(log *slog.Logger, writer *kafka.Writer) GoodsProducer {

	return GoodsProducer{
		writer: writer,
		log:    log,
	}
}

// А я точно не кринге снизу написал?
func CheckContextKafka(goodsproducer GoodsProducer, ctx context.Context) error {
	<-ctx.Done()

	err := goodsproducer.writer.Close()

	if err != nil {
		fmt.Errorf("не удалось закрыть writer: %w", err)
	}

	return nil
}

func (producer GoodsProducer) ProduceGoodsChanges(msgctx context.Context, message []int64) error {

	kafkaMessage := kafkaMessageType{
		goodsIds: message,
		dt:       time.Now(),
	}

	jsonmessage, err := json.Marshal(kafkaMessage)
	if err != nil {
		return err
	}

	err = producer.writer.WriteMessages(msgctx, kafka.Message{
		Key:   []byte(fmt.Sprintf("%s", time.Now())),
		Value: jsonmessage,
	})

	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Errorf("Контекст завершен по таймауту при отправке сообщения %g", err)
		} else {
			fmt.Errorf("Ошибка при отправке сообщения в Kafka: %v", err)
		}
	}

	return nil
}
