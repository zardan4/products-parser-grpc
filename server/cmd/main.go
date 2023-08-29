package main

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zardan4/products-parser-grpc/pkg/db"
	"github.com/zardan4/products-parser-grpc/server/internal/config"
	server_grpc "github.com/zardan4/products-parser-grpc/server/internal/server/grpc"
	"github.com/zardan4/products-parser-grpc/server/internal/service"
	storage "github.com/zardan4/products-parser-grpc/server/internal/storage/mongo"
)

func main() {
	cfg, err := config.New()

	if err != nil {
		logrus.Fatal(err.Error())
		return
	}

	dbClientWrapper, err := db.InitMongo(cfg.ConnLine)
	defer dbClientWrapper.Disconnect()
	if err != nil {
		logrus.Fatal(err.Error())
		return
	}

	db := dbClientWrapper.Client.Database(cfg.DBName)

	storage := storage.NewStorage(db)
	service := service.NewService(storage)
	server := server_grpc.NewServer(service)

	fmt.Printf("Server started at %s", time.Now().UTC())

	if err := server.ListenAndServe(cfg.Port); err != nil {
		logrus.Fatal(err)
	}
}
