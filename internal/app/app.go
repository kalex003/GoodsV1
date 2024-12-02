package app

import (
	"Goodsv1/internal/adapter/adapter"
	grpcapp "Goodsv1/internal/app/grpc"
	"Goodsv1/internal/kafka/goods_consumer"
	"Goodsv1/internal/kafka/goods_producer"
	"Goodsv1/internal/services/goods"
	"Goodsv1/internal/storage/postgres"
	kafkainit "Goodsv1/pkg/kafka"
	"github.com/segmentio/kafka-go"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App

	//тута рест сервис со сваггером
}

func New(
	grpcport int,
	log *slog.Logger,
	ConnString string, // чет сложно, тут надо разбираться до конца нормально
) (*App, *postgres.GoodsDb, *kafka.Reader, goods_consumer.GoodsConsumer, goods_producer.GoodsProducer) {

	GoodsDb, err := postgres.New(ConnString, log)

	if err != nil {
		panic(err)
	}

	Writer, err := kafkainit.NewKafkaWriter()
	if err != nil {
		panic(err)
	}

	Reader := kafkainit.NewKafkaReader()

	GoodsProducer := goods_producer.New(log, Writer)

	GoodsAdapter := adapter.New(GoodsProducer, log)

	GoodsConsumer := goods_consumer.New(GoodsAdapter, log)

	GoodsService := goods.New(log, GoodsDb, GoodsDb, GoodsDb, GoodsDb, GoodsAdapter) //надо адаптер добавить

	grpcApp := grpcapp.New(GoodsService, grpcport)

	return &App{
			GRPCServer: grpcApp,
		},
		GoodsDb,
		Reader,
		GoodsConsumer,
		GoodsProducer
}
