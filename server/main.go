package main

import (
	"context"

	"github.com/mikerjacobi/pingpen/server/pb"
	grpcj "github.com/zang-cloud/grpc-json"
)

type server struct{}

func (server *server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	resp := &pb.AddResponse{
		Sum: req.NumOne + req.NumTwo,
	}
	return resp, nil
}

func main() {
	grpcj.Serve(&server{}, grpcj.Port(":8082"))
}
