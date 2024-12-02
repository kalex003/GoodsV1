package entity

import (
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
