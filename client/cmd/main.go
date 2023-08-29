package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/zardan4/products-parser-grpc/pkg/core/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const ListOrFetch = "fetch"

// const ListOrFetch = "list"

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	godotenv.Load()

	conn, err := grpc.Dial(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatal(err)
		return
	}
	defer conn.Close()

	client := product.NewProductServiceClient(conn)

	if ListOrFetch == "fetch" {
		req := &product.FetchRequest{
			Url: "http://164.92.251.245:8080/api/v1/products/",
		}

		res, err := client.Fetch(context.TODO(), req)
		if err != nil {
			logrus.Fatal(err)
			return
		}

		fmt.Printf("%d products has been added, %d products has been updated. %d products has been deleted.\nTOTALLY PRODUCTS IN DB: %d",
			res.AddsAmount,
			res.UpdatesAmount,
			res.DeletesAmount,
			res.TotalAmount)
	}

	if ListOrFetch == "list" {
		// sorting mode choose
		mode, err := product.ToPbSortMode(product.SORT_MODE_LAST_CHANGE)
		if err != nil {
			logrus.Fatal(err)
			return
		}

		req := &product.ListRequest{
			Reversed:   true,
			Mode:       mode,
			PageSize:   0,
			PageOffset: 0,
		}

		res, err := client.List(context.TODO(), req)
		if err != nil {
			logrus.Fatal(err)
			return
		}

		for _, item := range res.Items {
			fmt.Printf("[NAME] %s; [PRICE] %d; [CHANGES] %d; [LAST CHANGE]: %s\n",
				item.GetName(), item.GetPrice(), item.GetChangeAmount(), item.GetLastChange().AsTime().Format(time.ANSIC))
		}
	}
}
