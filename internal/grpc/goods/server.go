package grpcgoods

import (
	goods1 "Goodsv1/grpc/gen/go/goods.v1"
	"Goodsv1/internal/entity"
	"Goodsv1/internal/grpc/models"
	"Goodsv1/internal/services/interfaces"
	"context"
	"errors"
	"google.golang.org/grpc"
)

type serverAPI struct { //для реализация интерфейсов, сгенерированнхы прото файлом
	goods1.UnimplementedGoodsServer
	goods interfaces.GoodsService // serverAPI зависит от Goods
}

func Register(gRPCServer *grpc.Server, goods interfaces.GoodsService) {
	goods1.RegisterGoodsServer(gRPCServer, &serverAPI{goods: goods})
}

func (s *serverAPI) Insert(ctx context.Context, request *goods1.InsertRequest) (*goods1.InsertResponse, error) {

	var err error
	var response []int64
	err = ValidateInsert(request)

	if err != nil {
		return nil, err
	}

	response, err = s.goods.InsertGoods(ctx, models.ConvertInsertRequestToGoodsSlice(request))

	if err != nil {
		return nil, err
	}

	return models.ConvertGoodsToInsertResponse(response), nil
}

func ValidateInsert(request *goods1.InsertRequest) error {
	if request == nil {
		return errors.New("request is nil")
	}
	return nil
}

func (s *serverAPI) Update(ctx context.Context, request *goods1.UpdateRequest) (*goods1.UpdateResponse, error) {

	var err error
	err = ValidateUpdate(request)

	if err != nil {
		return nil, err
	}

	err = s.goods.UpdateGoods(ctx, models.ConvertUpdateRequestToGoodsSlice(request))

	if err != nil {
		return nil, err
	}

	return &goods1.UpdateResponse{}, nil
}

func ValidateUpdate(request *goods1.UpdateRequest) error {
	if request == nil {
		return errors.New("request is nil")
	}
	return nil
}

func (s *serverAPI) GetById(ctx context.Context, request *goods1.GetByIdRequest) (*goods1.GetResponse, error) {

	var err error
	var response []entity.Good
	err = ValidateGetById(request)

	if err != nil {
		return nil, err
	}

	response, err = s.goods.GetGoodsByIds(ctx, request.GetGoodsId())

	if err != nil {
		return nil, err
	}

	return models.ConvertGoodsSliceToGetResponse(response), nil
}

func ValidateGetById(request *goods1.GetByIdRequest) error {
	if request == nil {
		return errors.New("request is nil")
	}
	return nil
}

func (s *serverAPI) GetByPlace(ctx context.Context, request *goods1.GetByPlaceRequest) (*goods1.GetResponse, error) {

	var err error
	var response []entity.Good
	err = ValidateGetByPlace(request)

	if err != nil {
		return nil, err
	}

	response, err = s.goods.GetGoodsByPlace(ctx, request.GetPlaceId())

	if err != nil {
		return nil, err
	}

	return models.ConvertGoodsSliceToGetResponse(response), nil
}

func ValidateGetByPlace(request *goods1.GetByPlaceRequest) error {
	if request == nil {
		return errors.New("request is nil")
	}
	return nil
}

func (s *serverAPI) GetByTare(ctx context.Context, request *goods1.GetByTareRequest) (*goods1.GetResponse, error) {

	var err error
	var response []entity.Good
	err = ValidateGetByTare(request)

	if err != nil {
		return nil, err
	}

	response, err = s.goods.GetGoodsByTare(ctx, request.GetTareId())

	if err != nil {
		return nil, err
	}

	return models.ConvertGoodsSliceToGetResponse(response), nil
}

func ValidateGetByTare(request *goods1.GetByTareRequest) error {
	if request == nil {
		return errors.New("request is nil")
	}
	return nil
}

func (s *serverAPI) GetHistory(ctx context.Context, request *goods1.GetHistoryRequest) (*goods1.GetResponse, error) {

	var err error
	var response []entity.Good
	err = ValidateGetHistory(request)

	if err != nil {
		return nil, err
	}

	response, err = s.goods.GetGoodsHistory(ctx, request.GetGoodsId())

	if err != nil {
		return nil, err
	}

	return models.ConvertGoodsSliceToGetResponse(response), nil
}

func ValidateGetHistory(request *goods1.GetHistoryRequest) error {
	if request == nil {
		return errors.New("request is nil")
	}
	return nil
}

func (s *serverAPI) Delete(ctx context.Context, request *goods1.DeleteRequest) (*goods1.DeleteResponse, error) {

	var err error
	err = ValidateDelete(request)

	if err != nil {
		return nil, err
	}

	err = s.goods.DeleteGoods(ctx, models.ConvertDeleteRequestToGoodsSlice(request))

	if err != nil {
		return nil, err
	}

	return &goods1.DeleteResponse{}, nil
}

func ValidateDelete(request *goods1.DeleteRequest) error {
	if request == nil {
		return errors.New("request is nil")
	}
	return nil
}
