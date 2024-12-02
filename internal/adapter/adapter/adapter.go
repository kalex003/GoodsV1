package adapter

import (
	adapter_models "Goodsv1/internal/adapter/models"
	"Goodsv1/internal/services/interfaces"
	"context"
	"log/slog"
)

type GoodsProducer interface { // интерфейс, который должен будет выполнять exporter
	ProduceGoodsChanges(context.Context, []int64) error
}

type GoodsAdapter struct {
	log           *slog.Logger
	GoodsProducer GoodsProducer
	GoodsService  interfaces.GoodsService // dependency injection
}

func New(producer GoodsProducer, log *slog.Logger) *GoodsAdapter { //мб стоит чекнуть, что указатель не nil?

	return &GoodsAdapter{GoodsProducer: producer,
		log: log}
}

func (ga *GoodsAdapter) ExportGoodsChanges(ctx context.Context, changes []int64) error {

	//AdapterChanges := adapter_models.ConvertSliceGoodsToAdaptermodel(changes)
	err := ga.GoodsProducer.ProduceGoodsChanges(ctx, changes)

	if err != nil {
		return err
	}

	return nil

}

func (ga *GoodsAdapter) ImportGoodsChanges(ctx context.Context, changes []adapter_models.Good) error {

	err := ga.GoodsService.ImportGoodsChanges(ctx, adapter_models.ConvertAdaptermodelToSliceGoods(changes))

	if err != nil {
		return err
	}
	return nil

}
