package todoserviceserver

import (
	"context"
	"database/sql"
	"log"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/awakair/awakair_todo_bot/api/todo-service"
)

type TodoServiceServer struct {
	Db *sql.DB
	pb.UnimplementedTodoServiceServer
}

type DbQuery = string

const (
	InsertNewUser DbQuery = "INSERT INTO users (id, language_code, utc_offset) VALUES ($1::bigint, $2, $3);"
)

func (s *TodoServiceServer) CreateUser(ctx context.Context, to_create *pb.User) (*emptypb.Empty, error) {
	v, err := protovalidate.New()

	if err != nil {
		log.Printf("failed to initialize validator at CreateUser: %v", err)

		return nil, status.Error(codes.Unknown, err.Error())
	}

	if err = v.Validate(to_create); err != nil {
		log.Printf("validation failed at CreateUser: %v", err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = s.Db.ExecContext(
		ctx, string(InsertNewUser), to_create.GetId(), to_create.GetLanguageCode(), to_create.GetUtcOffset(),
	)

	if err != nil {
		log.Printf(
			"Failed to create new user with id %v, language code %v and utc_offset %v: %v",
			to_create.GetId(), to_create.GetLanguageCode(), to_create.GetUtcOffset(),
			err,
		)
		return nil, status.Error(codes.Unknown, err.Error())
	}

	log.Printf(
		"Created new user with id %v, language code %v and utc_offset %v",
		to_create.GetId(), to_create.GetLanguageCode(), to_create.GetUtcOffset(),
	)

	return nil, nil
}

func (s *TodoServiceServer) EditUserSettings(ctx context.Context, in *pb.User) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditUserSettings not implemented")
}
func (s *TodoServiceServer) AddReminder(ctx context.Context, in *pb.Reminder) (*pb.ReminderId, error) {

	return nil, status.Errorf(codes.Unimplemented, "method AddReminder not implemented")
}
func (s *TodoServiceServer) RemoveReminder(ctx context.Context, in *pb.ReminderId) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveReminder not implemented")
}
func (s *TodoServiceServer) GetRemindersByUserId(in *pb.UserId, in1 pb.TodoService_GetRemindersByUserIdServer) error {
	return status.Errorf(codes.Unimplemented, "method GetRemindersByUserId not implemented")
}
