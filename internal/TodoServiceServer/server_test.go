package todoserviceserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pb "github.com/awakair/awakair_todo_bot/api/todo-service"
)

type StubRepo struct {
	SetUserFunc func(context.Context, *pb.User) error
}

func (sr *StubRepo) SetUser(ctx context.Context, user *pb.User) error {
	return sr.SetUserFunc(ctx, user)
}

func server(ctx context.Context, sr *StubRepo) (pb.TodoServiceClient, func()) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	baseServer := grpc.NewServer()
	pb.RegisterTodoServiceServer(baseServer, New(sr))
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

	return client, closer
}

func TestTodoServiceServer_SetUser(t *testing.T) {
	ctx := context.Background()

	usersCount := 0

	sr := &StubRepo{SetUserFunc: func(context.Context, *pb.User) error {
		usersCount++

		return nil
	}}

	client, closer := server(ctx, sr)
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
			_, err := client.SetUser(ctx, user)

			if status.Code(err) != codes.InvalidArgument {
				t.Errorf("expected InvalidArgument with user %+v got %v", user, err)
			}
		}

		backupUsersCount := usersCount
		usersCount = 0
		if backupUsersCount != 0 {

			t.Errorf("Did not expected calls of repo.SetUser, got %v calls", backupUsersCount)
		}
	})

	t.Run("regular user", func(t *testing.T) {
		user := &pb.User{
			Id:           0,
			LanguageCode: wrapperspb.String("en"),
			UtcOffset:    wrapperspb.Int32(0),
		}

		_, err := client.SetUser(ctx, user)

		if err != nil {
			t.Errorf("did not expect error with user %+v got %v", user, err)
		}

		backupUsersCount := usersCount
		usersCount = 0
		if backupUsersCount != 1 {

			t.Errorf("Expected 1 call of repo.SetUser, got %v calls", backupUsersCount)
		}
	})

	t.Run("partial filled user", func(t *testing.T) {
		users := []*pb.User{
			{Id: 0},
			{Id: 0, LanguageCode: wrapperspb.String("en")},
			{Id: 0, UtcOffset: wrapperspb.Int32(0)},
		}

		for _, user := range users {
			_, err := client.SetUser(ctx, user)

			if err != nil {
				t.Errorf("did not expect error with user %+v got %v", user, err)
			}
		}

		backupUsersCount := usersCount
		usersCount = 0
		if backupUsersCount != len(users) {

			t.Errorf("Expected %v calls of repo.SetUser, got %v calls", len(users), backupUsersCount)
		}
	})

	t.Run("db returns error", func(t *testing.T) {
		user := &pb.User{
			Id:           0,
			LanguageCode: wrapperspb.String("en"),
			UtcOffset:    wrapperspb.Int32(0),
		}

		sr.SetUserFunc = func(context.Context, *pb.User) error {
			usersCount++

			return fmt.Errorf("oops...")
		}

		_, err := client.SetUser(ctx, user)

		if status.Code(err) != codes.ResourceExhausted {
			t.Errorf("expected ResourceExhausted error with user %+v got %v", user, err)
		}

		backupUsersCount := usersCount
		usersCount = 0
		if backupUsersCount != 1 {
			t.Errorf("Expected 1 call of repo.SetUser, got %v calls", backupUsersCount)
		}
	})
}
