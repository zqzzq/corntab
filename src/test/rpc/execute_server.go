package main

import (
	"test/rpc/pb"
	"context"
	"net"
	"log"
	"google.golang.org/grpc"
	"fmt"
)

const (
	port = ":50053"
)


type ExecuteServiceServer struct {
}

func (ess *ExecuteServiceServer) Execute(ctx context.Context, req *pb.ExecuteRequest) (*pb.ExecuteResponse, error) {

	fmt.Println("hello world")
	return &pb.ExecuteResponse{Result:"success"}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterExecuteServiceServer(s, &ExecuteServiceServer{})
	if err = s.Serve(lis); err != nil{
		log.Fatalf("failed to serve: %v", err)
	}
}