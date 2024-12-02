package goods

import (
	"Goodsv1/internal/entity"
	"context"
	"log/slog"
)

type GoodsService struct {
	log           *slog.Logger
	GoodsInserter GoodsInserter //storage
	GoodsUpdater  GoodsUpdater  //storage
	GoodsGetter   GoodsGetter   //storage
	GoodsDeleter  GoodsDeleter  //storage
	GoodsExporter GoodsExporter //adapter
}

type GoodsInserter interface {
	InsertGoods(context.Context, []entity.Good) ([]int64, error)
}

type GoodsUpdater interface {
	GoodsUpdate(context.Context, []entity.Good) error
}
type GoodsGetter interface {
	GetGoodsByIds(context.Context, []int64) ([]entity.Good, error)
	GetGoodsByPlace(context.Context, int64) ([]entity.Good, error)
	GetGoodsByTare(context.Context, int64) ([]entity.Good, error)
	GetGoodsHistory(context.Context, int64) ([]entity.Good, error)
}

type GoodsDeleter interface {
	DeleteGoods(context.Context, []entity.Good) error
}

type GoodsExporter interface {
	ExportGoodsChanges(context.Context, []int64) error
}

func New(log *slog.Logger, GoodsInserter GoodsInserter, GoodsUpdater GoodsUpdater, GoodsGetter GoodsGetter, GoodsDeleter GoodsDeleter, GoodsExporter GoodsExporter) *GoodsService {
	return &GoodsService{
		log:           log,
		GoodsInserter: GoodsInserter,
		GoodsUpdater:  GoodsUpdater,
		GoodsGetter:   GoodsGetter,
		GoodsDeleter:  GoodsDeleter,
		GoodsExporter: GoodsExporter,
	}
}

func (s *GoodsService) InsertGoods(ctx context.Context, goods []entity.Good) ([]int64, error) {

	goodsIds, err := s.GoodsInserter.InsertGoods(ctx, goods)
	if err != nil {
		return nil, err
	}

	err = s.GoodsExporter.ExportGoodsChanges(ctx, goodsIds)
	if err != nil {
		return nil, err
	}

	return goodsIds, nil
}

func (s *GoodsService) UpdateGoods(ctx context.Context, goods []entity.Good) error {

	err := s.GoodsUpdater.GoodsUpdate(ctx, goods)
	if err != nil {
		return err
	}

	return nil
}

func (s *GoodsService) DeleteGoods(ctx context.Context, goods []entity.Good) error {

	err := s.GoodsDeleter.DeleteGoods(ctx, goods)
	if err != nil {
		return err
	}

	return nil
}

func (s *GoodsService) GetGoodsByIds(ctx context.Context, goodsIds []int64) ([]entity.Good, error) {

	goods, err := s.GoodsGetter.GetGoodsByIds(ctx, goodsIds)
	if err != nil {
		return nil, err
	}

	return goods, nil
}

func (s *GoodsService) GetGoodsByPlace(ctx context.Context, placeId int64) ([]entity.Good, error) {

	goods, err := s.GoodsGetter.GetGoodsByPlace(ctx, placeId)
	if err != nil {
		return nil, err
	}

	return goods, nil
}

func (s *GoodsService) GetGoodsByTare(ctx context.Context, tareId int64) ([]entity.Good, error) {

	goods, err := s.GoodsGetter.GetGoodsByTare(ctx, tareId)
	if err != nil {
		return nil, err
	}

	return goods, nil
}

func (s *GoodsService) GetGoodsHistory(ctx context.Context, goodsId int64) ([]entity.Good, error) {

	goods, err := s.GoodsGetter.GetGoodsHistory(ctx, goodsId)
	if err != nil {
		return nil, err
	}

	return goods, nil
}

func (s *GoodsService) ExportGoodsChanges(ctx context.Context, goods []int64) error {

	err := s.GoodsExporter.ExportGoodsChanges(ctx, goods)
	if err != nil {
		return err
	}

	return nil
}

func (s *GoodsService) ImportGoodsChanges(ctx context.Context, goods []entity.Good) error {

	if len(goods) != 0 {
		s.log.Info(("Сообщение от кафка обработано"))
	}

	return nil

	//надо подумать что ещё можно тут сделать, после того как приняли сообщение от самого же себя (кринге правда, но ладн)
}
