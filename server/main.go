package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	pb "github.com/Fahad-I-Khan/grpc-user-service/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedUserServiceServer
	store map[string]user
	mu    sync.RWMutex
}

func newServer() *server {
	return &server{
		store: make(map[string]user),
	}
}

type user struct {
	Id    string
	Name  string
	Email string
}

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	log.Println("Received request for ID:", req.Id)

	s.mu.RLock()
	defer s.mu.RUnlock()

	// For demonstration, we return a static user. In a real application, you'd query a database.
	u, exists := s.store[req.Id]
	if !exists {
		// return nil, grpc.Errorf(grpc.Code(pb.ErrorCode_USER_NOT_FOUND), "User not found")
		log.Println("User with ID not found:", req.Id)
		return nil, status.Error(codes.NotFound, "User not found")
	}

	return &pb.UserResponse{
		Id:    u.Id,
		Name:  u.Name,
		Email: u.Email,
	}, nil
}

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	log.Println("Received request to create user with name:", req.Name)

	s.mu.Lock()
	defer s.mu.Unlock()

	// For demonstration, we generate a simple ID. In a real application, you'd use a proper ID generator.
	id := uuid.New().String()
	s.store[id] = user{
		Id:    id,
		Name:  req.Name,
		Email: req.Email,
	}
	return &pb.UserResponse{
		Id:    id,
		Name:  req.Name,
		Email: req.Email,
	}, nil
}

func (s *server) ListUsers(ctx context.Context, _ *pb.Empty) (*pb.ListUsersResponse, error) {
	log.Println("Received request to list users")
	
	s.mu.RLock()
	defer s.mu.RUnlock()

	var users []*pb.UserResponse
	for _, u := range s.store {
		users = append(users, &pb.UserResponse{
			Id:    u.Id,
			Name:  u.Name,
			Email: u.Email,
		})
	}

	return &pb.ListUsersResponse{Users: users}, nil
}

func (s *server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.Empty, error) {
	log.Println("Received request to delete user with ID:", req.Id)
	
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.store[req.Id]; !exists {
		log.Println("User with ID not found:", req.Id)
		return nil, status.Error(codes.NotFound, "User not found")
	}

	delete(s.store, req.Id)
	
	return &pb.Empty{}, nil
}


func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, newServer())

	reflection.Register(grpcServer) // This enables grpcurl to work

	go func() {
		log.Println("gRPC server running on port 50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("Server stopped")
}


