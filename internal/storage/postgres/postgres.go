package postgres

import (
	"Goodsv1/internal/entity"
	"Goodsv1/internal/storage/models"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"strings"
)

type GoodsDb struct {
	log *slog.Logger
	Db  *pgxpool.Pool // используем pgxpool для пула подключений
}

func New(conn string, log *slog.Logger) (*GoodsDb, error) {

	pool, err := pgxpool.New(context.Background(), conn)
	if err != nil {
		return nil, err
	}

	return &GoodsDb{Db: pool, log: log}, nil
}

func (db *GoodsDb) LogGoodsChange(ctx context.Context, goods []models.Good) error {

	ins_stat := squirrel.Insert("goods.goodslog AS g").
		Columns("goods_id", "place_id", "employee_id", "tare_id", "dt", "is_del").
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING g.log_dt;")

	for _, i := range goods {
		ins_stat = ins_stat.Values(i.GoodsId, i.PlaceId, i.EmployeeId, i.TareId, i.Dt, i.IsDel)
	}

	sql, args, err := ins_stat.ToSql()
	if err != nil {
		return fmt.Errorf("LogGoodsChange: unable to build query: %w", err)
	}

	// Выполняем запрос
	rows, err := db.Db.Query(ctx, sql, args...)
	if err != nil {
		fmt.Errorf("LogGoodsChange: %w", err)
	}

	defer rows.Close()

	var ins_count int

	for rows.Next() {
		ins_count++
	}

	if ins_count != len(goods) {
		fmt.Errorf("LogGoodsChange: %w", "Ошибка логирования")
	}

	return nil

}

func (db *GoodsDb) InsertGoods(ctx context.Context, goods []entity.Good) ([]int64, error) {

	ins_stat := squirrel.Insert("goods.goods AS g").
		Columns("place_id", "employee_id", "tare_id").
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING g.goods_id")

	for _, i := range goods {
		ins_stat = ins_stat.Values(i.PlaceId, i.EmployeeId, i.TareId)
	}

	sql, args, err := ins_stat.ToSql()
	if err != nil {
		return []int64{}, fmt.Errorf("InsertGoods: unable to build query: %w", err)
	}

	// Выполняем запрос
	rows, err := db.Db.Query(ctx, sql, args...)
	if err != nil {
		fmt.Errorf("InsertGoods: %w", err)
	}

	defer rows.Close()

	var insertedGoods []models.Good

	for rows.Next() {
		var currentRow models.Good
		if err := rows.Scan(&currentRow.GoodsId); err != nil {
			return []int64{}, err
		}

		insertedGoods = append(insertedGoods, currentRow)
	}

	err = db.LogGoodsChange(ctx, insertedGoods)

	if err != nil {
		return []int64{}, err
	}

	var goodsIds []int64

	for _, row := range insertedGoods {
		goodsIds = append(goodsIds, row.GoodsId)
	}

	return goodsIds, nil
}

func (db *GoodsDb) GoodsUpdate(ctx context.Context, goods []entity.Good) error {

	/*		upd_stat, _, err := squirrel.StatementBuilder.Update("goods.goods AS g").
				SetMap("dt", squirrel.Expr("CURRENT_TIMESTAMP AT TIME ZONE 'Europe/Moscow'").
				Set("place_id", squirrel.Expr("c.place_id")).
				Set("ch_employee_id", squirrel.Expr("c.ch_employee_id")).
				Set("tare_id", squirrel.Expr("c.tare_id")).
				Set("tare_type", squirrel.Expr("c.tare_type")).
				Ta
			Where("g.goods_id = c.goods_id").
				PlaceholderFormat(squirrel.Dollar).
				ToSql()

			if err != nil {
				return fmt.Errorf("GoodsUpdate: unable to build query: %w", err)
			}*/

	var values []string
	for _, g := range models.ConvertSliceGoodsToDbModel(goods) {
		values = append(values, fmt.Sprintf("(%d, %d, %d, %d)", g.GoodsId, g.PlaceId, g.EmployeeId, *g.TareId))
	}

	valuesSQL := strings.Join(values, ",")

	query := fmt.Sprintf(`
    UPDATE goods.goods AS g
    SET dt = CURRENT_TIMESTAMP AT TIME ZONE 'Europe/Moscow',
        place_id = c.place_id,
        employee_id = c.employee_id,
        tare_id = c.tare_id
    FROM (VALUES %s) AS c(goods_id, place_id, employee_id, tare_id)
    WHERE g.goods_id = c.goods_id
    RETURNING g.*;`, valuesSQL)

	db.log.Info(fmt.Sprintf("Попытка обновить, запрос: %s", query))
	rows, err := db.Db.Query(ctx, query)

	if err != nil {
		fmt.Errorf("GoodsUpdate: %w", err)
	}

	defer rows.Close()

	var updatedGoods []models.Good

	for rows.Next() {
		var currentRow models.Good
		if err := rows.Scan(&currentRow.GoodsId, &currentRow.PlaceId, &currentRow.EmployeeId, &currentRow.Dt, &currentRow.TareId, &currentRow.IsDel); err != nil {
			return err
		}

		updatedGoods = append(updatedGoods, currentRow)
	}

	err = db.LogGoodsChange(ctx, updatedGoods)

	if err != nil {
		return err
	}

	return nil
}

func (db *GoodsDb) DeleteGoods(ctx context.Context, goods []entity.Good) error {

	/*			upd_stat, _, err := squirrel.StatementBuilder.Update("goods.goods AS g").
					SetMap(squirrel.Eq{"is_del": goodsDb., "y": 2}).
					Where("g.goods_id = c.goods_id").
					PlaceholderFormat(squirrel.Dollar).
					ToSql()

				if err != nil {
					return fmt.Errorf("GoodsUpdate: unable to build query: %w", err)
				}
	*/

	var values []string
	for _, g := range models.ConvertSliceGoodsToDbModel(goods) {
		values = append(values, fmt.Sprintf("(%d, %t)", g.GoodsId, g.IsDel))
	}

	valuesSQL := strings.Join(values, ",")

	query := fmt.Sprintf(`
    UPDATE goods.goods AS g
    SET dt = CURRENT_TIMESTAMP AT TIME ZONE 'Europe/Moscow',
        is_del = c.is_del
    FROM (VALUES %s) AS c(goods_id, is_del)
    WHERE g.goods_id = c.goods_id
    RETURNING g.*;`, valuesSQL)

	db.log.Info(fmt.Sprintf("Попытка обновить признак удаления, запрос: %s", query))

	rows, err := db.Db.Query(ctx, query)

	if err != nil {
		fmt.Errorf("GoodsDelete: %w", err)
	}

	defer rows.Close()

	var updatedGoods []models.Good

	for rows.Next() {
		var currentRow models.Good
		if err := rows.Scan(&currentRow.GoodsId, &currentRow.PlaceId, &currentRow.EmployeeId, &currentRow.Dt, &currentRow.TareId, &currentRow.IsDel); err != nil {
			return err
		}

		updatedGoods = append(updatedGoods, currentRow)
	}

	err = db.LogGoodsChange(ctx, updatedGoods)

	if err != nil {
		return err
	}

	return nil
}

func (db *GoodsDb) GetGoodsByIds(ctx context.Context, goodsIds []int64) ([]entity.Good, error) {

	var err error
	var sql string
	var rows pgx.Rows
	if len(goodsIds) == 0 {
		sql, _, err = squirrel.StatementBuilder.
			Select("g.goods_id", "g.place_id", "g.employee_id", "g.tare_id", "g.dt", "g.is_del").
			From("goods.goods AS g").
			ToSql()

		if err != nil {
			return []entity.Good{}, fmt.Errorf("SelectGoodsByIds: unable to build query: %w", err)
		}

		rows, err = db.Db.Query(ctx, sql)

	} else {
		sql, _, err = squirrel.StatementBuilder.
			Select("g.goods_id", "g.place_id", "g.employee_id", "g.tare_id", "g.dt", "g.is_del").
			From("goods.goods AS g").
			Where("g.goods_id = ANY(?)").
			PlaceholderFormat(squirrel.Dollar).
			ToSql()

		if err != nil {
			return []entity.Good{}, fmt.Errorf("SelectGoodsByIds: unable to build query: %w", err)
		}

		rows, err = db.Db.Query(ctx, sql, goodsIds)
	}

	if err != nil {
		return []entity.Good{}, fmt.Errorf("SelectGoodsByIds: %w", err)
	}
	defer rows.Close()

	var result []entity.Good

	for rows.Next() {

		var row entity.Good

		if err := rows.Scan(&row.GoodsId, &row.PlaceId, &row.EmployeeId, &row.TareId, &row.Dt, &row.IsDel); err != nil {
			return []entity.Good{}, err
		}
		result = append(result, row)
	}
	return result, nil

}

func (db *GoodsDb) GetGoodsByPlace(ctx context.Context, placeId int64) ([]entity.Good, error) {

	sql, args, err := squirrel.StatementBuilder.
		Select("g.goods_id", "g.place_id", "g.employee_id", "g.tare_id", "g.dt", "g.is_del").
		From("goods.goods AS g").
		Where("g.place_id = ?", placeId).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return []entity.Good{}, fmt.Errorf("GetGoodsByPlace: unable to build query: %w", err)
	}

	rows, err := db.Db.Query(ctx, sql, args...)
	if err != nil {
		return []entity.Good{}, fmt.Errorf("GetGoodsByPlace: %w", err)
	}
	defer rows.Close()

	var result []entity.Good

	for rows.Next() {
		var row entity.Good
		if err := rows.Scan(&row.GoodsId, &row.PlaceId, &row.EmployeeId, &row.TareId, &row.Dt, &row.IsDel); err != nil {
			return []entity.Good{}, err
		}
		result = append(result, row)
	}
	return result, nil

}

func (db *GoodsDb) GetGoodsByTare(ctx context.Context, tareId int64) ([]entity.Good, error) {

	sql, args, err := squirrel.StatementBuilder.
		Select("g.goods_id", "g.place_id", "g.employee_id", "g.tare_id", "g.dt", "g.is_del").
		From("goods.goods AS g").
		Where("g.place_id = ?", tareId).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return []entity.Good{}, fmt.Errorf("GetGoodsByTare: unable to build query: %w", err)
	}

	rows, err := db.Db.Query(ctx, sql, args...)
	if err != nil {
		return []entity.Good{}, fmt.Errorf("GetGoodsByTare: %w", err)
	}
	defer rows.Close()

	var result []entity.Good

	for rows.Next() {
		var row entity.Good
		if err := rows.Scan(&row.GoodsId, &row.PlaceId, &row.EmployeeId, &row.TareId, &row.Dt, &row.IsDel); err != nil {
			return []entity.Good{}, err
		}
		result = append(result, row)
	}
	return result, nil

}

func (db *GoodsDb) GetGoodsHistory(ctx context.Context, goodsId int64) ([]entity.Good, error) {

	sql, args, err := squirrel.StatementBuilder.
		Select("g.goods_id", "g.place_id", "g.employee_id", "g.tare_id", "g.dt", "g.is_del").
		From("goods.goodslog AS g").
		Where("g.goods_id = ?", goodsId).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return []entity.Good{}, fmt.Errorf("GetGoodsHistory: unable to build query: %w", err)
	}

	rows, err := db.Db.Query(ctx, sql, args...)
	if err != nil {
		return []entity.Good{}, fmt.Errorf("GetGoodsHistory: %w", err)
	}
	defer rows.Close()

	var result []entity.Good

	for rows.Next() {
		var row entity.Good
		if err := rows.Scan(&row.GoodsId, &row.PlaceId, &row.EmployeeId, &row.TareId, &row.Dt, &row.IsDel); err != nil {
			return []entity.Good{}, err
		}
		result = append(result, row)
	}
	return result, nil

}
