package server_grpc

import (
	"context"

	"github.com/zardan4/products-parser-grpc/pkg/core/product"
	"github.com/zardan4/products-parser-grpc/server/internal/service"
)

type ProductsServer struct {
	// product.UnimplementedProductServiceServer
	product.ProductServiceServer

	service *service.Service
}

func NewProductsServer(service *service.Service) *ProductsServer {
	return &ProductsServer{
		service: service,
	}
}

func (s *ProductsServer) Fetch(ctx context.Context, req *product.FetchRequest) (*product.FetchResponse, error) {
	return s.service.Product.Fetch(ctx, req)
}

func (s *ProductsServer) List(ctx context.Context, req *product.ListRequest) (*product.ListResponse, error) {
	return s.service.Product.List(ctx, req)
}
