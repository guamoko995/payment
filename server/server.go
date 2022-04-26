package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"

	proto "payment/proto"
	db "payment/server/database"
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
func (s *PaymentServer) UpSum(ctx context.Context, req *proto.UpRequest) (*proto.Enum, error) {
	err := s.Storage.UpBalance(req.ID, int64(req.Sum))
	if err != nil {
		return &proto.Enum{}, err
	}

	return &proto.Enum{}, nil
}

func (s *PaymentServer) SumTransfer(ctx context.Context, req *proto.TransferRequest) (*proto.Enum, error) {
	err := s.Storage.SumTransfer(req.SenderID, req.GeterID, int64(req.Sum))
	if err != nil {
		return &proto.Enum{}, err
	}

	return &proto.Enum{}, nil
}

// создает экземпляр gRPC прокси-сервиса
func newServer() *PaymentServer {
	b, err := db.NewDB()
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
	proto.RegisterPaymentServer(grpcServer, newServer())
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := grpcServer.Serve(lis)
		if err != nil {
			log.Printf("failed to listen: %v\n", err)
		}
		wg.Done()
	}()
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt, syscall.SIGTERM)
	log.Printf("got signal %v, attempting graceful stop\n", <-interruptCh)
	grpcServer.GracefulStop()
	wg.Wait()
	log.Println("Server gracefully stopped")
}
