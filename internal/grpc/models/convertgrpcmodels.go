package models

import (
	goods1 "Goodsv1/grpc/gen/go/goods.v1"
	"Goodsv1/internal/entity"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertOneInsertRequestToGood(req *goods1.OneInsertRequest) entity.Good {
	return entity.Good{
		PlaceId:    req.PlaceId,
		EmployeeId: req.EmployeeId,
		TareId:     req.TareId,
	}
}

func ConvertInsertRequestToGoodsSlice(reqs *goods1.InsertRequest) []entity.Good {
	var goods []entity.Good
	for _, req := range reqs.GetStructs() {
		good := ConvertOneInsertRequestToGood(req)
		goods = append(goods, good)
	}

	return goods
}

func ConvertOneUpdateRequestToGood(req *goods1.OneUpdateRequest) entity.Good {
	return entity.Good{
		GoodsId:    req.GoodsId,
		PlaceId:    req.PlaceId,
		EmployeeId: req.EmployeeId,
		TareId:     req.TareId,
	}
}

func ConvertUpdateRequestToGoodsSlice(reqs *goods1.UpdateRequest) []entity.Good {
	var goods []entity.Good
	for _, req := range reqs.GetStructs() {
		good := ConvertOneUpdateRequestToGood(req)
		goods = append(goods, good)
	}

	return goods
}

func ConvertOneDeleteRequestToGood(req *goods1.OneDeleteRequest) entity.Good {
	return entity.Good{
		GoodsId: req.GoodsId,
		IsDel:   req.IsDel,
	}
}

func ConvertDeleteRequestToGoodsSlice(reqs *goods1.DeleteRequest) []entity.Good {
	var goods []entity.Good
	for _, req := range reqs.GetStructs() {
		good := ConvertOneDeleteRequestToGood(req)
		goods = append(goods, good)
	}

	return goods
}

func ConvertGoodsToInsertResponse(resps []int64) *goods1.InsertResponse {
	var goodsIds []int64
	for _, resp := range resps {
		goodId := resp
		goodsIds = append(goodsIds, goodId)
	}

	return &goods1.InsertResponse{
		GoodsId: goodsIds,
		Dt:      timestamppb.Now(),
	}
}

func ConvertGoodsToOneGetResponse(resp entity.Good) *goods1.OneGetResponse {
	return &goods1.OneGetResponse{
		GoodsId:    resp.GoodsId,
		PlaceId:    resp.PlaceId,
		EmployeeId: resp.EmployeeId,
		TareId:     resp.TareId,
		Dt:         timestamppb.New(resp.Dt),
		IsDel:      resp.IsDel,
	}
}

func ConvertGoodsSliceToGetResponse(resps []entity.Good) *goods1.GetResponse {
	var goods []*goods1.OneGetResponse
	for _, resp := range resps {
		good := ConvertGoodsToOneGetResponse(resp)
		goods = append(goods, good)
	}

	return &goods1.GetResponse{
		Structs: goods,
	}
}
