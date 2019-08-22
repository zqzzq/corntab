package main

import (

	"context"
	"net"
	"log"
	"google.golang.org/grpc"
	"fmt"
	"corntab/src/samples/pb"
)

const (
	port = ":50052"
)


type ExecuteServiceServer struct {
}

func (ess *ExecuteServiceServer) Execute(ctx context.Context, req *pb.ExecuteRequest) (*pb.ExecuteResponse, error) {

	fmt.Println("hello ", req.Params)
	return &pb.ExecuteResponse{Output:"success"},nil
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