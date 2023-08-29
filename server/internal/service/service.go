package service

import (
	"context"

	"github.com/zardan4/products-parser-grpc/pkg/core/product"
	storage "github.com/zardan4/products-parser-grpc/server/internal/storage/mongo"
)

type Product interface {
	Fetch(ctx context.Context, req *product.FetchRequest) (*product.FetchResponse, error)
	List(ctx context.Context, req *product.ListRequest) (*product.ListResponse, error)
}

type Service struct {
	Product
}

func NewService(storage *storage.Storage) *Service {
	return &Service{
		Product: NewProductService(storage),
	}
}
