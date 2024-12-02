package goods_consumer

import (
	"Goodsv1/internal/adapter/adapter"
	"Goodsv1/internal/kafka/kafka_models"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log/slog"
	"time"
)

var msgch chan kafka.Message
var msgctx context.Context

type GoodsConsumer struct {
	log *slog.Logger
	//reader  *kafka.Reader
	adapter *adapter.GoodsAdapter //dependency injection
}

func New( /*reader *kafka.Reader, */ adapter *adapter.GoodsAdapter, log *slog.Logger) GoodsConsumer {

	return GoodsConsumer{
		//	reader:  reader,
		adapter: adapter,
		log:     log,
	}
}

// Worker — основной процесс чтения сообщений из Kafka и записи их в базу данных
func ConsumerWorker(ctx context.Context, log *slog.Logger, reader *kafka.Reader) {

	defer func() {
		if err := reader.Close(); err != nil {
			log.Error("Ошибка при закрытии читателя Kafka", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Info("Kafka consumer остановлен по сигналу контекста")
			return
		default:

			msg, err := reader.ReadMessage(ctx) //ридер убрать что ли надо из структура консьюмер? Сначала попробую метод написать правильно, а потом уже интерфейс под него
			if err != nil {
				if err == context.DeadlineExceeded {
					log.Error("Таймаут при чтении сообщения из Kafka")
				} else {
					log.Error("Ошибка при чтении сообщения из Kafka: %v", err)
				}
				continue
			}

			log.Info("Получено сообщение из Kafka: %s", msg.Value)

			msgctx, _ = context.WithTimeout(ctx, 30*time.Second) //хотелось бы как-нить контекст создать чтобы не все время обрабатывалоь

			msgch <- msg

			/*cancel()*/
		}
	}
}

func (gc GoodsConsumer) ConsumeMessage(ctx context.Context) {

	var importedgoods []kafka_models.Good

	for {
		select {
		case <-ctx.Done():
			gc.log.Info("Kafka consumer остановлен по сигналу контекста")
		case <-msgch:

			msg := <-msgch

			err := json.Unmarshal(msg.Value, &importedgoods)
			if err != nil {
				gc.log.Info("ошибка при десериализации сообщения: %w", err)

			}

			err = gc.adapter.ImportGoodsChanges(msgctx, kafka_models.ConvertKafkamodelToSliceGoods(importedgoods))

			if err != nil {
				gc.log.Info("ошибка при обработке сообщения в слое адаптера: %w", err)

			}

		}
	}

}
