package service

import (
	"context"
	 "hunger4data/pb/subcription"
	"subscription/internal/adapters/external"
	"subscription/internal/adapters/model"
	"subscription/internal/adapters/repo"

	"github.com/google/uuid"
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
	.UnimplementedSubscription_ServiceServer
	RDB *redis.Client
}

func NewSubService(repo repo.SubscriptionRepo, RDB *redis.Client) SubService {
	return SubService{repo: repo, RDB: RDB}
}

func (s *SubService) Get_Countries(ctx context.Context, req *.Empty) (*.Get_Countries_Response, error) {
	data, err := external.GetHumDataRedis(s.RDB)
	if err != nil {
		if err.Error() == "cache empty" {
			return &.Get_Countries_Response{
				Message: "Data is currently being synchronized, please try again shortly.",
			}, nil
		}

	}
	return data, nil
}

func (s *SubService) Create(ctx context.Context, req *.Subscription_Request) (*.Subscription_Response, error) {
	if req.IdCountry == "" && req.IdUser == "" {
		return &.Subscription_Response{Message: "user id and country id  are required"}, status.Error(codes.InvalidArgument, "user id and country id  are required")
	}
	iduser, _ := uuid.Parse(req.IdUser)
	idcountry, _ := uuid.Parse(req.IdCountry)
	subs := model.Subscription{
		Id_user:    iduser,
		Id_country: idcountry,
	}
	err := s.repo.CreateSubcription(subs)
	if err != nil {
		return &.Subscription_Response{Message: "Error Creating subscription"}, status.Error(codes.Internal, "Error Creating subscription")
	}
	return &.Subscription_Response{Message: "user id and country id  are required"}, nil

}

func (s *SubService) GetByID(ctx context.Context, req *.Subscription_Request) (*.Get_Subscription_BY_ID_Response, error) {
	if req.IdUser == "" {
		return &.Get_Subscription_BY_ID_Response{Message: "user id is required"}, status.Error(codes.InvalidArgument, "user id is required")
	}
	iduser, _ := uuid.Parse(req.IdUser)
	data, err := s.repo.GetBySubscriptionUserID(iduser)
	if err != nil {
		return &.Get_Subscription_BY_ID_Response{Message: "Error Creating subscription"}, status.Error(codes.Internal, "Error Creating subscription")
	}
	protoSubs := []*.Subscription{}
	for _, sub := range data {
		protoSubs = append(protoSubs, &.Subscription{
			Id:        sub.Id.String(),
			IdUser:    sub.Id_user.String(),
			IdCountry: sub.Id_country.String(),
		})
	}
	return &.Get_Subscription_BY_ID_Response{Subscription: protoSubs, Message: "Succes Creating"}, nil

}

func (s *SubService) Update(ctx context.Context, req *.Subscription_Request) (*.Subscription_Response, error) {
	if req.Id == "" {
		return &.Subscription_Response{Message: "id is required"}, status.Error(codes.InvalidArgument, "id is required")
	}
	iduser, _ := uuid.Parse(req.IdUser)
	idcountry, _ := uuid.Parse(req.IdCountry)
	subs := model.Subscription{
		Id_user:    iduser,
		Id_country: idcountry,
	}
	id, _ := uuid.Parse(req.Id)
	err := s.repo.UpdateSubscription(id, subs)
	if err != nil {
		return &.Subscription_Response{Message: "Error Updating subscription"}, status.Error(codes.Internal, "Error Updating subscription")
	}
	return &.Subscription_Response{Message: "Success Updating"}, nil

}

func (s *SubService) Delete(ctx context.Context, req *.Subscription_Request) (*.Subscription_Response, error) {
	if req.Id == "" {
		return &.Subscription_Response{Message: "id is required"}, status.Error(codes.InvalidArgument, "id is required")
	}
	id, _ := uuid.Parse(req.Id)
	err := s.repo.DeleteSubscription(id)
	if err != nil {
		return &.Subscription_Response{Message: "Error Deleting subscription"}, status.Error(codes.Internal, "Error Deleting subscription")
	}
	return &.Subscription_Response{Message: "succes deleting"}, nil

}
