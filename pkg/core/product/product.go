package product

import (
	"errors"
	"time"
)

const (
	SORT_MODE_NAME           = "NAME"
	SORT_MODE_PRICE          = "PRICE"
	SORT_MODE_CHANGES_AMOUNT = "CHANGES_AMOUNT"
	SORT_MODE_LAST_CHANGE    = "LAST_CHANGE"
)

var (
	sortModes = map[string]ListRequest_SortingMode{
		SORT_MODE_NAME:           ListRequest_NAME,
		SORT_MODE_PRICE:          ListRequest_PRICE,
		SORT_MODE_CHANGES_AMOUNT: ListRequest_CHANGES_AMOUNT,
		SORT_MODE_LAST_CHANGE:    ListRequest_LAST_CHANGE,
	}
)

type ProductItem struct {
	Name          string    `bson:"name"`
	Price         int64     `bson:"price"`
	ChangesAmount int64     `bson:"changes_amount"`
	LastChange    time.Time `bson:"last_change"`
}

func ToPbSortMode(mode string) (ListRequest_SortingMode, error) {
	val, ex := sortModes[mode]
	if !ex {
		return 0, errors.New("invalid sort mode")
	}

	return val, nil
}
