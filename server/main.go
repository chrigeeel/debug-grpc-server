package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/chrigeeel/debug-grpc-server/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type server struct {
	pb.UnimplementedIPServiceServer
}

func (s *server) GetIP(ctx context.Context, req *pb.IPRequest) (*pb.IPResponse, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to get peer from context")
	}

	clientIP, _, err := net.SplitHostPort(p.Addr.String())
	if err != nil {
		return nil, fmt.Errorf("failed to parse client IP: %v", err)
	}

	return &pb.IPResponse{Ip: clientIP}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":29381")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterIPServiceServer(grpcServer, &server{})
	log.Println("Server is running on port :29381")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
