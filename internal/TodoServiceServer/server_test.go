package todoserviceserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"regexp"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	pb "github.com/awakair/awakair_todo_bot/api/todo-service"
	dbqueries "github.com/awakair/awakair_todo_bot/internal/DbQueries"
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

	t.Run("wrong users", func(t *testing.T) {
		users := []*pb.User{
			{},
			{
				Id:           0,
				LanguageCode: "hello world",
				UtcOffset:    0,
			},
			{
				Id:           0,
				LanguageCode: "en",
				UtcOffset:    2109,
			},
		}

		for _, user := range users {
			_, err := client.CreateUser(ctx, user)

			if status.Code(err) != codes.InvalidArgument {
				t.Errorf("expected InvalidArgument while inserting %v got %v", user, err)
			}
		}
	})

	t.Run("regular user", func(t *testing.T) {
		user := &pb.User{
			Id:           0,
			LanguageCode: "en",
			UtcOffset:    0,
		}

		mock.ExpectExec(regexp.QuoteMeta(dbqueries.InsertNewUser)).WithArgs(
			user.GetId(), user.GetLanguageCode(), user.GetUtcOffset(),
		).WillReturnResult(
			sqlmock.NewResult(user.GetId(), 1),
		)

		_, err := client.CreateUser(ctx, user)

		if err != nil {
			t.Errorf("did not expect error while inserting user %v got %v", user, err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("mock expectations weren't met, error: %v", err)
		}
	})

	t.Run("regular user", func(t *testing.T) {
		user := &pb.User{
			Id:           0,
			LanguageCode: "en",
			UtcOffset:    0,
		}

		mock.ExpectExec(regexp.QuoteMeta(dbqueries.InsertNewUser)).WithArgs(
			user.GetId(), user.GetLanguageCode(), user.GetUtcOffset(),
		).WillReturnError(
			fmt.Errorf("oops..."),
		)

		_, err := client.CreateUser(ctx, user)

		if status.Code(err) != codes.ResourceExhausted {
			t.Errorf("expected ResourceExhausted error while inserting user %v got %v", user, err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("mock expectations weren't met, error: %v", err)
		}
	})
}
