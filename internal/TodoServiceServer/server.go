package todoserviceserver

import (
	"context"
	"log"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/awakair/awakair_todo_bot/api/todo-service"
)

type Repo interface {
	SetUser(context.Context, *pb.User) error
	// GetUser(context.Context, int) error
}

type TodoServiceServer struct {
	repo Repo
	pb.UnimplementedTodoServiceServer
}

func New(repo Repo) *TodoServiceServer {
	return &TodoServiceServer{repo: repo}
}

func (s *TodoServiceServer) SetUser(ctx context.Context, user *pb.User) (_ *emptypb.Empty, err error) {
	defer func() {
		if err != nil {
			log.Printf("Error in SetUser with user %+v: %v", user, err)
		} else {
			log.Printf("SetUser with user %+v was successful", user)
		}
	}()

	v, err := protovalidate.New()

	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	if err = v.Validate(user); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.repo.SetUser(ctx, user)

	if err != nil {
		return nil, status.Error(codes.ResourceExhausted, err.Error())
	}

	return nil, nil
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
