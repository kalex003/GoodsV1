package adapter_interfaces

import (
	adapter_models "Goodsv1/internal/adapter/models"
	"context"
)

type GoodsConsumer interface {
	ImportGoodsChanges(context.Context, []adapter_models.Good) error
}
