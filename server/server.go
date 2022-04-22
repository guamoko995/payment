package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	proto "payment_service/payment"
	db "payment_service/server/database"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
	port       = flag.Int("port", 50051, "The server port")
	database   = flag.String("database", "", "The file of database SQLite")
)

// gRPC прокси-сервис для осуществления пополнения
// аккаунтов и переводов.
type PaymentServer struct {
	proto.UnimplementedPaymentServer
	Storage *db.DB
}

// Реализация соотвествующих RPC сервиса.
func (s *PaymentServer) UpAccount(ctx context.Context, req *proto.UpRequest) (*proto.Enum, error) {
	err := s.Storage.UpBalance(req.ID, int64(req.Sum))
	if err != nil {
		return &proto.Enum{}, err
	}

	return &proto.Enum{}, nil
}

func (s *PaymentServer) AmountTransfer(ctx context.Context, req *proto.TransferRequest) (*proto.Enum, error) {
	err := s.Storage.AmountTransfer(req.SenderID, req.GeterID, int64(req.Sum))
	if err != nil {
		return &proto.Enum{}, err
	}

	return &proto.Enum{}, nil
}

// создает экземпляр gRPC прокси-сервиса
func newServer(database string) *PaymentServer {
	b, err := db.NewDB(database)
	if err != nil {
		panic(err)
	}
	s := &PaymentServer{Storage: b}
	return s
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterPaymentServer(grpcServer, newServer(*database))
	grpcServer.Serve(lis)
}
