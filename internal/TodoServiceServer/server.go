package todoserviceserver

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/awakair/awakair_todo_bot/api/todo-service"
	dbqueries "github.com/awakair/awakair_todo_bot/internal/DbQueries"
)

type TodoServiceServer struct {
	Db *sql.DB
	pb.UnimplementedTodoServiceServer
}

func (s *TodoServiceServer) DbActionWithUser(
	ctx context.Context, user *pb.User,
	operationName string, userFilledQuery dbqueries.UserFilled,
) (_ *emptypb.Empty, err error) {
	defer func() {
		if err != nil {
			log.Printf("Error in %v with user %+v: %v", operationName, user, err)
		} else {
			log.Printf("Operation %v with user %+v was successful", operationName, user)
		}
	}()

	v, err := protovalidate.New()

	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	if err = v.Validate(user); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var queryResult sql.Result

	if user.GetLanguageCode() == nil && user.GetUtcOffset() == nil {
		if userFilledQuery.Empty == "" {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("can't use %v without fields", operationName))
		}
		queryResult, err = s.Db.ExecContext(
			ctx, userFilledQuery.Empty,
			user.GetId(),
		)
	} else if user.GetLanguageCode() != nil && user.GetUtcOffset() != nil {
		if userFilledQuery.Full == "" {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("can't use %v with all fields", operationName))
		}
		queryResult, err = s.Db.ExecContext(
			ctx, userFilledQuery.Full,
			user.GetId(),
			user.GetLanguageCode().GetValue(), user.GetUtcOffset().GetValue(),
		)
	} else if user.GetLanguageCode() != nil {
		if userFilledQuery.LanguageCode == "" {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("can't use %v only with LanguageCode", operationName))
		}
		queryResult, err = s.Db.ExecContext(
			ctx, userFilledQuery.LanguageCode,
			user.GetId(),
			user.GetLanguageCode().GetValue(),
		)
	} else {
		if userFilledQuery.UtcOffset == "" {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("can't use %v only with UtcOffset", operationName))
		}
		queryResult, err = s.Db.ExecContext(
			ctx, userFilledQuery.UtcOffset,
			user.GetId(),
			user.GetUtcOffset().GetValue(),
		)
	}

	if err != nil {
		return nil, status.Error(codes.ResourceExhausted, err.Error())
	}

	rowsAffected, err := queryResult.RowsAffected()

	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	if rowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return nil, nil
}

func (s *TodoServiceServer) CreateUser(ctx context.Context, to_create *pb.User) (*emptypb.Empty, error) {
	return s.DbActionWithUser(ctx, to_create, "TodoServiceServer.CreateUser", dbqueries.InsertUser)
}

func (s *TodoServiceServer) EditUserSettings(ctx context.Context, to_update *pb.User) (*emptypb.Empty, error) {
	return s.DbActionWithUser(ctx, to_update, "TodoServiceServer.UpdateUser", dbqueries.UpdateUser)
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
