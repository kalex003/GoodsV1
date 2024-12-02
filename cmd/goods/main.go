package main

import (
	"Goodsv1/config"
	"Goodsv1/internal/app"
	"Goodsv1/internal/kafka/goods_consumer"
	"Goodsv1/internal/kafka/goods_producer"
	"Goodsv1/internal/storage/partition_maker"
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

var ConnString = MustGetEnv("DATABASE_URL")

func MustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Missing required environment variable: %s", key)
	}

	return value
}

// Константы для Kafka

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	log := setupLogger(cfg.Env)

	log.Info("Starting application",
		slog.String("env", cfg.Env),
		slog.Any("env", cfg),
		slog.Int("port", cfg.GRPC.Port),
	)

	application, db, reader, consumer, producer := app.New(cfg.GRPC.Port, log, ConnString) //добавил для крончика

	go application.GRPCServer.MustRun() //в параллельном режиме от остальной программы обрабатываем запросы

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT) //обработка заверщаюего сигнала (одного из двух сигтерм или сигинт)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go partition_maker.Worker(ctx, log, db)

	go goods_consumer.ConsumerWorker(ctx, log, reader)

	go consumer.ConsumeMessage(ctx)

	go goods_producer.CheckContextKafka(producer, ctx)
	<-stop //пока в этот канал что-то не придет (сигнал о зщавершении программы), мы тут просто будем висеть и ждать, а сверху будет работать го рутина

	application.GRPCServer.Stop()

}

func setupLogger(env string) *slog.Logger { //сложно чот
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
