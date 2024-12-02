package kafka_models

import (
	adapter_models "Goodsv1/internal/adapter/models"
	"time"
)

type Good struct {
	GoodsId    int64
	PlaceId    int64
	EmployeeId int64
	TareId     *int64
	Dt         time.Time
	IsDel      bool
}

func ConvertGoodsToKafkamodel(good adapter_models.Good) Good {
	return Good{
		GoodsId:    good.GoodsId,
		PlaceId:    good.PlaceId,
		Dt:         good.Dt,
		EmployeeId: good.EmployeeId,
		TareId:     good.TareId,
		IsDel:      good.IsDel,
	}
}

func ConvertSliceGoodsToKafkamodel(goods []adapter_models.Good) []Good {
	var result []Good
	for _, info := range goods {
		good := ConvertGoodsToKafkamodel(info)
		result = append(result, good)
	}
	return result
}

func ConvertKafkamodelToGoods(good Good) adapter_models.Good {
	return adapter_models.Good{
		GoodsId:    good.GoodsId,
		PlaceId:    good.PlaceId,
		Dt:         good.Dt,
		EmployeeId: good.EmployeeId,
		TareId:     good.TareId,
		IsDel:      good.IsDel,
	}
}

func ConvertKafkamodelToSliceGoods(goods []Good) []adapter_models.Good {
	var result []adapter_models.Good
	for _, info := range goods {
		good := ConvertKafkamodelToGoods(info)
		result = append(result, good)
	}
	return result
}
