package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mikerjacobi/pingpen/server/pb"
	"github.com/sirupsen/logrus"
	grpcj "github.com/zang-cloud/grpc-json"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
}

type server struct{}

func (s *server) Add(ctx context.Context, req *pb.AddRequest) (*pb.MathResponse, error) {
	resp := &pb.MathResponse{
		Sum: req.NumOne + req.NumTwo,
	}
	fmt.Println("ADD")
	return resp, nil
}

func (s *server) Sub(ctx context.Context, req *pb.SubRequest) (*pb.MathResponse, error) {
	resp := &pb.MathResponse{
		Sum: req.NumOne - req.NumTwo,
	}
	fmt.Println("SUB")
	return resp, nil
}

func main() {
	s := &server{}
	function := os.Getenv("function")
	switch function {
	case "add":
		lambda.Start(s.Add)
	case "sub":
		lambda.Start(s.Sub)
	default:
		grpcj.Serve(s, grpcj.Port(":8082"))
	}
}
