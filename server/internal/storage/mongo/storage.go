package storage

import (
	"context"

	"github.com/zardan4/products-parser-grpc/pkg/core/product"
	"go.mongodb.org/mongo-driver/mongo"
)

type Product interface {
	FetchProductsByURL(url string) ([]product.ProductItem, error)
	List(ctx context.Context, opts ListParams) ([]product.ProductItem, error)

	UpdateByNames(ctx context.Context, products []product.ProductItem) (UpdatingResponse, error)
	FindAllToUpdateAndNew(ctx context.Context, products []product.ProductItem) (FindAllToUpdateAndNewOutput, error)
}

type Storage struct {
	Product
}

func NewStorage(db *mongo.Database) *Storage {
	return &Storage{
		Product: NewProductStorage(db),
	}
}
