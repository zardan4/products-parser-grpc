package storage

import (
	"context"
	"encoding/csv"
	"net/http"
	"strconv"
	"time"

	collections "github.com/zardan4/products-parser-grpc/pkg/core/dbCollections"
	"github.com/zardan4/products-parser-grpc/pkg/core/product"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductStorage struct {
	db *mongo.Database
}

func NewProductStorage(db *mongo.Database) *ProductStorage {
	return &ProductStorage{
		db: db,
	}
}

type FindAllToUpdateAndNewOutput struct {
	// total
	New     []product.ProductItem // all items there are not in db
	Old     []product.ProductItem // all items there are in db but not synchronized
	OldSync []product.ProductItem // all items there are in db and synchronized

	Deleted []product.ProductItem // items that have been deleted after last sync
}

// returns all products that should be updated(has different price in db) and products that are new in db
func (s *ProductStorage) FindAllToUpdateAndNew(ctx context.Context, products []product.ProductItem) (FindAllToUpdateAndNewOutput, error) {
	coll := s.db.Collection(collections.ProductsCollection, nil)

	cur, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return FindAllToUpdateAndNewOutput{}, err
	}
	defer cur.Close(ctx)

	var allProdsDB []product.ProductItem

	err = cur.All(ctx, &allProdsDB)
	if err != nil {
		return FindAllToUpdateAndNewOutput{}, err
	}

	var res FindAllToUpdateAndNewOutput

	// filter total items(New, Old, OldSync)
	for _, p := range products {
		found := false

		for _, p2 := range allProdsDB {
			// if old item
			if p2.Name == p.Name {
				found = true

				if p2.Price != p.Price {
					// if old not sync item
					res.Old = append(res.Old, p)
				} else {
					// if old sync item
					res.OldSync = append(res.OldSync, p)
				}

				break
			}
		}

		if !found {
			// if new item
			res.New = append(res.New, p)
		}
	}

	// filter deleted items
	for _, pDb := range allProdsDB {
		found := false

		for _, p := range products {
			if p.Name == pDb.Name {
				found = true
				break
			}
		}

		if !found {
			res.Deleted = append(res.Deleted, pDb)
		}
	}

	return res, nil
}

type UpdatingResponse struct {
	UpdAmount     int64
	AddAmount     int64
	DeletedAmount int64
	TotalAmount   int64
}

// if had anything to update - update. if item is new - add
func (s *ProductStorage) UpdateByNames(ctx context.Context, products []product.ProductItem) (UpdatingResponse, error) {
	coll := s.db.Collection(collections.ProductsCollection, nil)

	splitted, err := s.FindAllToUpdateAndNew(ctx, products)
	if err != nil {
		return UpdatingResponse{}, err
	}

	// update products
	var updModels []mongo.WriteModel

	for _, p := range splitted.Old {
		filter := bson.M{"name": p.Name}
		upd := bson.M{
			"$set": bson.M{
				"price":       p.Price,
				"last_change": time.Now(),
			},
			"$inc": bson.M{"changes_amount": 1},
		}

		updModel := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(upd).
			SetUpsert(false)

		updModels = append(updModels, updModel)
	}

	if len(updModels) != 0 {
		_, err = coll.BulkWrite(ctx, updModels)
		if err != nil {
			return UpdatingResponse{}, err
		}
	}

	// add new products
	var addModels []mongo.WriteModel

	for _, p := range splitted.New {
		filter := bson.M{"name": p.Name}
		upd := bson.M{
			"$set": bson.M{
				"name":           p.Name,
				"price":          p.Price,
				"changes_amount": 1,
				"last_change":    time.Now(),
			},
		}

		updModel := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(upd).
			SetUpsert(true)

		addModels = append(addModels, updModel)
	}

	if len(addModels) != 0 {
		_, err = coll.BulkWrite(ctx, addModels)
		if err != nil {
			return UpdatingResponse{}, err
		}
	}

	// delete products
	var deleteFilters []string

	for _, p := range splitted.Deleted {
		deleteFilters = append(deleteFilters, p.Name)
	}

	if len(splitted.Deleted) > 0 {
		_, err = coll.DeleteMany(ctx, bson.M{
			"name": bson.M{"$in": deleteFilters},
		})
		if err != nil {
			return UpdatingResponse{}, err
		}
	}

	return UpdatingResponse{
		UpdAmount:     int64(len(updModels)),
		AddAmount:     int64(len(addModels)),
		DeletedAmount: int64(len(splitted.Deleted)),
		TotalAmount:   int64(len(splitted.New) + len(splitted.Old) + len(splitted.OldSync)),
	}, nil
}

func (s *ProductStorage) FetchProductsByURL(url string) ([]product.ProductItem, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	reader := csv.NewReader(res.Body)
	reader.Comma = ';' // custom separator

	productsUnformatted, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var products []product.ProductItem

	for _, record := range productsUnformatted {
		price, _ := strconv.Atoi(record[1])
		priceInt64 := int64(price)

		var product = product.ProductItem{
			Name:  record[0],
			Price: priceInt64,
		}

		products = append(products, product)
	}

	return products, nil
}

type ListParams struct {
	Reversed   bool
	Mode       string
	PageSize   int64
	PageOffset int64
}

func (s *ProductStorage) List(ctx context.Context, opts ListParams) ([]product.ProductItem, error) {
	coll := s.db.Collection(collections.ProductsCollection)

	var findOptions []*options.FindOptions

	// setting reversed/normal search
	var reversedModeNum int
	if !opts.Reversed {
		reversedModeNum = 1
	} else {
		reversedModeNum = -1
	}

	// setting sort method
	var sortModeStruct bson.D

	switch opts.Mode {
	case product.SORT_MODE_NAME:
		sortModeStruct = bson.D{{"name", reversedModeNum}}
	case product.SORT_MODE_PRICE:
		sortModeStruct = bson.D{{"price", reversedModeNum}}
	case product.SORT_MODE_CHANGES_AMOUNT:
		sortModeStruct = bson.D{{"changes_amount", reversedModeNum}}
	case product.SORT_MODE_LAST_CHANGE:
		sortModeStruct = bson.D{{"last_change", reversedModeNum}}
	default:
		sortModeStruct = bson.D{{"name", reversedModeNum}}
	}

	findOptions = append(findOptions, options.Find().SetSort(sortModeStruct))

	// pagination settings
	// offset from start
	findOptions = append(findOptions, options.Find().SetSkip(opts.PageOffset))
	// limit return documents
	findOptions = append(findOptions, options.Find().SetLimit(opts.PageSize))

	cur, err := coll.Find(context.Background(), bson.M{}, findOptions...)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var allProdsDB []product.ProductItem

	err = cur.All(context.Background(), &allProdsDB)
	if err != nil {
		return nil, err
	}

	return allProdsDB, nil
}
