package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	pb "go-test/proto"
)

var store = make(map[int64]map[string]string)

type Server struct {
	mx sync.RWMutex
}

func (s *Server) SendMessage(ctx context.Context, message *pb.Message) (*empty.Empty, error) {
	fmt.Println("SendMessage", message)
	s.mx.Lock()
	store[message.Index] = message.Data
	s.mx.Unlock()

	return &empty.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := Server{}

	grpcServer := grpc.NewServer()

	pb.RegisterServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
