package models

import (
	"Goodsv1/internal/entity"
	"time"
)

type Good struct {
	GoodsId    int64     `db:"goods_id"`
	PlaceId    int64     `db:"place_id"`
	EmployeeId int64     `db:"employee_id"`
	TareId     *int64    `db:"tare_id"`
	Dt         time.Time `db:"dt"`
	IsDel      bool      `db:"is_del"`
}

func ConvertGoodsToDbModel(good entity.Good) Good {
	return Good{
		GoodsId:    good.GoodsId,
		PlaceId:    good.PlaceId,
		Dt:         good.Dt,
		EmployeeId: good.EmployeeId,
		TareId:     good.TareId,
		IsDel:      good.IsDel,
	}
}

func ConvertSliceGoodsToDbModel(goods []entity.Good) []Good {
	var result []Good
	for _, info := range goods {
		good := ConvertGoodsToDbModel(info)
		result = append(result, good)
	}
	return result
}
