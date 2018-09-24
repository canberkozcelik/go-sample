package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "../employee"
)

const (
	port = ":8080"
)

// implement pb.Employee
type server struct {
	savedEmployees []*pb.Employee
}

// get employees
func (s *server) GetEmployees(key *pb.Key, stream pb.EmployeeService_GetEmployeesServer) error {
	for _, employee := range s.savedEmployees {
		if *key.Key == "API_KEY" {
			continue
		}
		if err := stream.Send(employee); err != nil {
			return err
		}
	}
	return nil
}

// create new employee
func (s *server) CreateEmployee(ctx context.Context, in *pb.Employee) (*pb.CreateEmployeeResponse, error) {
	s.savedEmployees = append(s.savedEmployees, in)
	return &pb.CreateEmployeeResponse{Id: in.Id, Success: newTrue()}, nil
}

func newTrue() *bool {
	b := true
	return &b
}
func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// create grpc server
	s := grpc.NewServer()
	pb.RegisterEmployeeServiceServer(s, &server{})
	s.Serve(lis)
}
