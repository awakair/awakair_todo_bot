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
	"google.golang.org/protobuf/types/known/wrapperspb"
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

	t.Run("wrong user", func(t *testing.T) {
		users := []*pb.User{
			{
				Id:           0,
				LanguageCode: wrapperspb.String("hello world"),
				UtcOffset:    wrapperspb.Int32(0),
			},
			{
				Id:           0,
				LanguageCode: wrapperspb.String("en"),
				UtcOffset:    wrapperspb.Int32(2109),
			},
			{
				Id:           0,
				LanguageCode: wrapperspb.String("russky"),
				UtcOffset:    nil,
			},
			{
				Id:           0,
				LanguageCode: nil,
				UtcOffset:    wrapperspb.Int32(314159),
			},
		}

		for _, user := range users {
			_, err := client.CreateUser(ctx, user)

			if status.Code(err) != codes.InvalidArgument {
				t.Errorf("expected InvalidArgument with user %+v got %v", user, err)
			}
		}
	})

	t.Run("regular user", func(t *testing.T) {
		user := &pb.User{
			Id:           0,
			LanguageCode: wrapperspb.String("en"),
			UtcOffset:    wrapperspb.Int32(0),
		}

		mock.ExpectExec(regexp.QuoteMeta(dbqueries.InsertUser.Full)).WithArgs(
			user.GetId(),
			user.GetLanguageCode().GetValue(), user.GetUtcOffset().GetValue(),
		).WillReturnResult(
			sqlmock.NewResult(user.GetId(), 1),
		)

		_, err := client.CreateUser(ctx, user)

		if err != nil {
			t.Errorf("did not expect error with user %+v got %v", user, err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("mock expectations weren't met with user %+v, error: %v", user, err)
		}
	})

	t.Run("partial filled user", func(t *testing.T) {
		users := []*pb.User{
			{Id: 0},
			{Id: 0, LanguageCode: wrapperspb.String("en")},
			{Id: 0, UtcOffset: wrapperspb.Int32(0)},
		}

		mock.ExpectExec(regexp.QuoteMeta(dbqueries.InsertUser.Empty)).WithArgs(
			users[0].GetId(),
		).WillReturnResult(
			sqlmock.NewResult(users[0].GetId(), 1),
		)

		mock.ExpectExec(regexp.QuoteMeta(dbqueries.InsertUser.LanguageCode)).WithArgs(
			users[1].GetId(),
			users[1].GetLanguageCode().GetValue(),
		).WillReturnResult(
			sqlmock.NewResult(users[1].GetId(), 1),
		)

		mock.ExpectExec(regexp.QuoteMeta(dbqueries.InsertUser.UtcOffset)).WithArgs(
			users[2].GetId(),
			users[2].GetUtcOffset().GetValue(),
		).WillReturnResult(
			sqlmock.NewResult(users[2].GetId(), 1),
		)

		for _, user := range users {
			_, err := client.CreateUser(ctx, user)

			if err != nil {
				t.Errorf("did not expect error with user %+v got %v", user, err)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("mock expectations weren't met, error: %v", err)
		}
	})

	t.Run("db returns error", func(t *testing.T) {
		user := &pb.User{
			Id:           0,
			LanguageCode: wrapperspb.String("en"),
			UtcOffset:    wrapperspb.Int32(0),
		}

		mock.ExpectExec(regexp.QuoteMeta(dbqueries.InsertUser.Full)).WithArgs(
			user.GetId(),
			user.GetLanguageCode().GetValue(), user.GetUtcOffset().GetValue(),
		).WillReturnError(
			fmt.Errorf("oops..."),
		)

		_, err := client.CreateUser(ctx, user)

		if status.Code(err) != codes.ResourceExhausted {
			t.Errorf("expected ResourceExhausted error with user %+v got %v", user, err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("mock expectations weren't met with user %+v, error: %v", user, err)
		}
	})
}

func TestTodoServiceServer_EditUserSettings(t *testing.T) {
	ctx := context.Background()

	client, closer, mock := server(ctx)
	defer closer()

	t.Run("wrong user", func(t *testing.T) {
		users := []*pb.User{
			{},
			{
				Id:           0,
				LanguageCode: wrapperspb.String("hello world"),
				UtcOffset:    wrapperspb.Int32(0),
			},
			{
				Id:           0,
				LanguageCode: wrapperspb.String("en"),
				UtcOffset:    wrapperspb.Int32(2109),
			},
			{
				Id:           0,
				LanguageCode: wrapperspb.String("russky"),
				UtcOffset:    nil,
			},
			{
				Id:           0,
				LanguageCode: nil,
				UtcOffset:    wrapperspb.Int32(314159),
			},
		}

		for _, user := range users {
			_, err := client.EditUserSettings(ctx, user)

			if status.Code(err) != codes.InvalidArgument {
				t.Errorf("expected InvalidArgument with user %+v got %v", user, err)
			}
		}
	})

	t.Run("regular user", func(t *testing.T) {
		user := &pb.User{
			Id:           0,
			LanguageCode: wrapperspb.String("en"),
			UtcOffset:    wrapperspb.Int32(0),
		}

		mock.ExpectExec(regexp.QuoteMeta(dbqueries.UpdateUser.Full)).WithArgs(
			user.GetId(),
			user.GetLanguageCode().GetValue(), user.GetUtcOffset().GetValue(),
		).WillReturnResult(
			sqlmock.NewResult(user.GetId(), 1),
		)

		_, err := client.EditUserSettings(ctx, user)

		if err != nil {
			t.Errorf("did not expect error with user %+v got %v", user, err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("mock expectations weren't met with user %+v, error: %v", user, err)
		}
	})

	t.Run("partial filled user", func(t *testing.T) {
		users := []*pb.User{
			{Id: 0, LanguageCode: wrapperspb.String("en")},
			{Id: 0, UtcOffset: wrapperspb.Int32(0)},
		}

		mock.ExpectExec(regexp.QuoteMeta(dbqueries.UpdateUser.LanguageCode)).WithArgs(
			users[0].GetId(),
			users[0].GetLanguageCode().GetValue(),
		).WillReturnResult(
			sqlmock.NewResult(users[0].GetId(), 1),
		)

		mock.ExpectExec(regexp.QuoteMeta(dbqueries.UpdateUser.UtcOffset)).WithArgs(
			users[1].GetId(),
			users[1].GetUtcOffset().GetValue(),
		).WillReturnResult(
			sqlmock.NewResult(users[1].GetId(), 1),
		)

		for _, user := range users {
			_, err := client.EditUserSettings(ctx, user)

			if err != nil {
				t.Errorf("did not expect error with user %+v got %v", user, err)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("mock expectations weren't met, error: %v", err)
		}
	})

	t.Run("db returns error", func(t *testing.T) {
		user := &pb.User{
			Id:           0,
			LanguageCode: wrapperspb.String("en"),
			UtcOffset:    wrapperspb.Int32(0),
		}

		mock.ExpectExec(regexp.QuoteMeta(dbqueries.UpdateUser.Full)).WithArgs(
			user.GetId(),
			user.GetLanguageCode().GetValue(), user.GetUtcOffset().GetValue(),
		).WillReturnError(
			fmt.Errorf("oops..."),
		)

		_, err := client.EditUserSettings(ctx, user)

		if status.Code(err) != codes.ResourceExhausted {
			t.Errorf("expected ResourceExhausted error with user %+v got %v", user, err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("mock expectations weren't met with user %+v, error: %v", user, err)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		user := &pb.User{
			Id:           0,
			LanguageCode: wrapperspb.String("en"),
			UtcOffset:    wrapperspb.Int32(0),
		}

		mock.ExpectExec(regexp.QuoteMeta(dbqueries.UpdateUser.Full)).WithArgs(
			user.GetId(),
			user.GetLanguageCode().GetValue(), user.GetUtcOffset().GetValue(),
		).WillReturnResult(
			sqlmock.NewResult(0, 0),
		)

		_, err := client.EditUserSettings(ctx, user)

		if status.Code(err) != codes.NotFound {
			t.Errorf("expected NotFound error with user %+v got %v", user, err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("mock expectations weren't met with user %+v, error: %v", user, err)
		}
	})
}
