package grpcapp

import (
	grpcgoods "Goodsv1/internal/grpc/goods"
	"Goodsv1/internal/services/interfaces"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type App struct {
	grpcServer *grpc.Server
	port       int
}

func New(goodsService interfaces.GoodsService, port int) *App {
	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)

	grpcgoods.Register(grpcServer, goodsService)
	return &App{
		grpcServer: grpcServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	addr, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := a.grpcServer.Serve(addr); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.grpcServer.GracefulStop()
}
