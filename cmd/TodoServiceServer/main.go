package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	pb "github.com/awakair/awakair_todo_bot/api/todo-service"
	"github.com/awakair/awakair_todo_bot/internal/postgresrepo"
	"github.com/awakair/awakair_todo_bot/internal/todoserviceserver"
)

var (
	port  = flag.Int("port", 50051, "The server port")
	dbUrl = "postgres://" +
		os.Getenv("DB_USER") +
		":" +
		os.Getenv("DB_PASSWORD") +
		"@localhost:5432/" +
		os.Getenv("DB_NAME")
)

func main() {
	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("cannot connect to database with url %s: %v", dbUrl, err)
	}
	defer pool.Close()

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTodoServiceServer(s, todoserviceserver.New(postgresrepo.New(pool)))

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
