package service

import (
	"context"
	"subscription/internal/adapters/external"
	"subscription/internal/adapters/model"
	"subscription/internal/adapters/repo"
	pb "subscription/proto/subcription"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// type SubscriptionServ interface {
// 	Create(ctx context.Context, req *pb.Subscription_Request) (*pb.Subscription_Response, error)
// 	GetByID(ctx context.Context, req *pb.Subscription_Request) (*pb.Get_Subscription_Response, error)
// 	Update(ctx context.Context, req *pb.Subscription_Request) (*pb.Subscription_Response, error)
// 	Delete(ctx context.Context, req *pb.Subscription_Request) (*pb.Subscription_Response, error)
// }

type SubService struct {
	repo repo.SubscriptionRepo
	pb.UnimplementedSubscription_ServiceServer
	RDB *redis.Client
}

func NewSubService(repo repo.SubscriptionRepo, RDB *redis.Client) SubService {
	return SubService{repo: repo, RDB: RDB}
}

func (s *SubService) Get_Countries(ctx context.Context, req *pb.Empty) (*pb.Get_Countries_Response, error) {
	data, err := external.GetHumDataRedis(s.RDB)
	if err != nil {
		if err.Error() == "cache empty" {
			return &pb.Get_Countries_Response{
				Message: "Data is currently being synchronized, please try again shortly.",
			}, nil
		}

	}
	return data, nil
}

func (s *SubService) Create(ctx context.Context, req *pb.Subscription_Request) (*pb.Subscription_Response, error) {
	if req.IdCountry == 0 && req.IdUser == 0 {
		return &pb.Subscription_Response{Message: "user id and country id  are required"}, status.Error(codes.InvalidArgument, "user id and country id  are required")
	}
	subs := model.Subscription{
		Id_user:    req.IdUser,
		Id_country: req.IdCountry,
	}
	err := s.repo.CreateSubcription(subs)
	if err != nil {
		return &pb.Subscription_Response{Message: "Error Creating subscription"}, status.Error(codes.Internal, "Error Creating subscription")
	}
	return &pb.Subscription_Response{Message: "user id and country id  are required"}, nil

}

func (s *SubService) GetByID(ctx context.Context, req *pb.Subscription_Request) (*pb.Get_Subscription_Response, error) {
	if req.IdUser == 0 {
		return &pb.Get_Subscription_Response{Message: "user id is required"}, status.Error(codes.InvalidArgument, "user id is required")
	}
	data, err := s.repo.GetBySubscriptionUserID(int(req.IdUser))
	if err != nil {
		return &pb.Get_Subscription_Response{Message: "Error Creating subscription"}, status.Error(codes.Internal, "Error Creating subscription")
	}

	return &pb.Get_Subscription_Response{Subscription: &pb.Subscription{Id: data.Id, IdUser: data.Id_user, IdCountry: data.Id_country}, Message: "Succes Creating"}, nil

}

func (s *SubService) Update(ctx context.Context, req *pb.Subscription_Request) (*pb.Subscription_Response, error) {
	if req.Id == 0 {
		return &pb.Subscription_Response{Message: "id is required"}, status.Error(codes.InvalidArgument, "id is required")
	}
	subs := model.Subscription{
		Id_user:    req.IdUser,
		Id_country: req.IdCountry,
	}
	err := s.repo.UpdateSubscription(int(req.Id), subs)
	if err != nil {
		return &pb.Subscription_Response{Message: "Error Updating subscription"}, status.Error(codes.Internal, "Error Updating subscription")
	}
	return &pb.Subscription_Response{Message: "Success Updating"}, nil

}

func (s *SubService) Delete(ctx context.Context, req *pb.Subscription_Request) (*pb.Subscription_Response, error) {
	if req.Id == 0 {
		return &pb.Subscription_Response{Message: "id is required"}, status.Error(codes.InvalidArgument, "id is required")
	}
	err := s.repo.DeleteSubscription(int(req.Id))
	if err != nil {
		return &pb.Subscription_Response{Message: "Error Deleting subscription"}, status.Error(codes.Internal, "Error Deleting subscription")
	}
	return &pb.Subscription_Response{Message: "succes deleting"}, nil

}
