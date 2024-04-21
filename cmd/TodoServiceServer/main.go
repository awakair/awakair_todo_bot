package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"

	pb "github.com/awakair/awakair_todo_bot/api/todo-service"
	todoserviceserver "github.com/awakair/awakair_todo_bot/internal/TodoServiceServer"
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
	db, err := sql.Open("pgx", dbUrl)
	if err != nil {
		log.Fatalf("Cannot connect to database with url %s: %v", dbUrl, err)
	}
	defer db.Close()

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTodoServiceServer(s, &todoserviceserver.TodoServiceServer{Db: db})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
