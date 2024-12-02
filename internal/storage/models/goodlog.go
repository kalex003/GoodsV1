package models

import "time"

type Goodlog struct {
	GoodsId    int64     `db:"goods_id"`
	PlaceId    int64     `db:"place_id"`
	EmployeeId int64     `db:"employee_id"`
	TareId     *int64    `db:"tare_id"`
	Dt         time.Time `db:"dt"`
	IsDel      bool      `db:"is_del"`
	logDt      time.Time `db:"log_dt"`
}
