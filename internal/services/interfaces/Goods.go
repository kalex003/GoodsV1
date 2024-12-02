package interfaces

import (
	"Goodsv1/internal/entity"
	"context"
)

type GoodsService interface {
	InsertGoods(context.Context, []entity.Good) ([]int64, error)
	UpdateGoods(context.Context, []entity.Good) error
	GetGoodsByIds(context.Context, []int64) ([]entity.Good, error)
	GetGoodsByPlace(context.Context, int64) ([]entity.Good, error)
	GetGoodsByTare(context.Context, int64) ([]entity.Good, error)
	GetGoodsHistory(context.Context, int64) ([]entity.Good, error)
	DeleteGoods(context.Context, []entity.Good) error
	ExportGoodsChanges(context.Context, []int64) error
	ImportGoodsChanges(context.Context, []entity.Good) error
}
