package interfaces

import (
	"context"
)

type GoodsConsumer interface {
	ConsumeMessage(context.Context) error
}
