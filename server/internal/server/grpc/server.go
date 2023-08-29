package server_grpc

import (
	"fmt"
	"net"

	"github.com/zardan4/products-parser-grpc/pkg/core/product"
	"github.com/zardan4/products-parser-grpc/server/internal/service"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server

	productServer *ProductsServer
}

func NewServer(service *service.Service) *Server {
	return &Server{
		grpcServer: grpc.NewServer(),

		productServer: NewProductsServer(service),
	}
}

func (s *Server) ListenAndServe(port int64) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	product.RegisterProductServiceServer(s.grpcServer, s.productServer)

	err = s.grpcServer.Serve(lis)

	return err
}
