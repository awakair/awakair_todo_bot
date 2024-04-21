package todoserviceserver

import (
	"context"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	pb "github.com/awakair/awakair_todo_bot/api/todo-service"
)

func server(ctx context.Context) (pb.TodoServiceClient, func(), sqlmock.Sqlmock) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("error creating db mock: %v", err)
	}

	baseServer := grpc.NewServer()
	pb.RegisterTodoServiceServer(baseServer, &TodoServiceServer{Db: db})
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Fatalf("error serving server: %v", err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error connecting to server: %v", err)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Fatalf("error closing listener: %v", err)
		}
		baseServer.Stop()
	}

	client := pb.NewTodoServiceClient(conn)

	return client, closer, mock
}

func TestTodoServiceServer_CreateUser(t *testing.T) {
	ctx := context.Background()

	client, closer, mock := server(ctx)
	defer closer()

	t.Run("empty user", func(t *testing.T) {
		_, err := client.CreateUser(ctx, &pb.User{})

		if status.Code(err) != codes.InvalidArgument {
			t.Errorf("expected InvalidArgument got %v", err)
		}
	})
}
