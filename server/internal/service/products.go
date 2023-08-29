package service

import (
	"context"

	"github.com/zardan4/products-parser-grpc/pkg/core/product"
	storage "github.com/zardan4/products-parser-grpc/server/internal/storage/mongo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProductService struct {
	storage *storage.Storage
}

func NewProductService(storage *storage.Storage) *ProductService {
	return &ProductService{storage: storage}
}

func (s *ProductService) Fetch(ctx context.Context, req *product.FetchRequest) (*product.FetchResponse, error) {
	products, err := s.storage.Product.FetchProductsByURL(req.GetUrl())
	if err != nil {
		return nil, err
	}

	res, err := s.storage.Product.UpdateByNames(ctx, products)

	return &product.FetchResponse{
		AddsAmount:    res.AddAmount,
		UpdatesAmount: res.UpdAmount,
		DeletesAmount: res.DeletedAmount,
		TotalAmount:   res.TotalAmount,
	}, err
}

func (s *ProductService) List(ctx context.Context, req *product.ListRequest) (*product.ListResponse, error) {
	reversedSort := req.GetReversed()
	sortMode := req.GetMode()
	pageSize := req.GetPageSize()
	pageOffset := req.GetPageOffset()

	items, err := s.storage.Product.List(ctx, storage.ListParams{
		Reversed:   reversedSort,
		Mode:       sortMode.String(),
		PageSize:   pageSize,
		PageOffset: pageOffset,
	})
	if err != nil {
		return nil, err
	}

	var resItems []*product.ListResponseItem

	for _, i := range items {
		resItems = append(resItems, &product.ListResponseItem{
			Name:         i.Name,
			Price:        i.Price,
			ChangeAmount: i.ChangesAmount,
			LastChange:   timestamppb.New(i.LastChange),
		})
	}

	return &product.ListResponse{
		Items: resItems,
	}, nil
}
