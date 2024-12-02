package adapter_models

import (
	"Goodsv1/internal/entity"
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

func ConvertGoodsToAdaptermodel(good entity.Good) Good {
	return Good{
		GoodsId:    good.GoodsId,
		PlaceId:    good.PlaceId,
		Dt:         good.Dt,
		EmployeeId: good.EmployeeId,
		TareId:     good.TareId,
		IsDel:      good.IsDel,
	}
}

func ConvertSliceGoodsToAdaptermodel(goods []entity.Good) []Good {
	var result []Good
	for _, info := range goods {
		good := ConvertGoodsToAdaptermodel(info)
		result = append(result, good)
	}
	return result
}

func ConvertAdaptermodelToGoods(good Good) entity.Good {
	return entity.Good{
		GoodsId:    good.GoodsId,
		PlaceId:    good.PlaceId,
		Dt:         good.Dt,
		EmployeeId: good.EmployeeId,
		TareId:     good.TareId,
		IsDel:      good.IsDel,
	}
}

func ConvertAdaptermodelToSliceGoods(goods []Good) []entity.Good {
	var result []entity.Good
	for _, info := range goods {
		good := ConvertAdaptermodelToGoods(info)
		result = append(result, good)
	}
	return result
}
