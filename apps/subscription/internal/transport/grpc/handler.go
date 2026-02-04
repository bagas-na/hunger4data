package grpcHandler

import (
	"context"
	"fmt"
	pb "hunger4data/pb/subcription"
	"subscription/internal/adapters/external"
	"subscription/internal/adapters/model"
	"subscription/internal/service"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SubService struct {
	serv service.SubServ
	pb.UnimplementedSubscription_ServiceServer
	RDB *redis.Client
}

func NewSubHand(serv service.SubServ, RDB *redis.Client) SubService {
	return SubService{serv: serv, RDB: RDB}
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
	iduser, _ := uuid.Parse(req.IdUser)
	idcountry, _ := uuid.Parse(req.IdCountry)
	subs := model.Subscription{
		Id_user:    iduser,
		Id_country: idcountry,
	}
	err := s.serv.CreateSubcription(subs)
	if err != nil {
		return &pb.Subscription_Response{Message: "Error Creating subscription"}, status.Error(codes.Internal, fmt.Sprintf("%s", err))
	}
	return &pb.Subscription_Response{Message: "user id and country id  are required"}, nil

}

func (s *SubService) GetByID(ctx context.Context, req *pb.Subscription_Request) (*pb.Get_Subscription_BY_ID_Response, error) {
	iduser, _ := uuid.Parse(req.IdUser)
	data, err := s.serv.GetSubscriptionByID(iduser)
	if err != nil {
		return &pb.Get_Subscription_BY_ID_Response{Message: "Error Getting subscription"}, status.Error(codes.Internal, fmt.Sprintf("%s", err))
	}
	protoSubs := []*pb.Subscription{}
	for _, sub := range data {
		protoSubs = append(protoSubs, &pb.Subscription{
			Id:        sub.Id.String(),
			IdUser:    sub.Id_user.String(),
			IdCountry: sub.Id_country.String(),
		})
	}
	return &pb.Get_Subscription_BY_ID_Response{Subscription: protoSubs, Message: "Success Getting Subscription"}, nil

}

func (s *SubService) Update(ctx context.Context, req *pb.Subscription_Request) (*pb.Subscription_Response, error) {
	idcountry, _ := uuid.Parse(req.IdCountry)
	subs := model.Subscription{
		Id_country: idcountry,
	}
	id, _ := uuid.Parse(req.Id)
	err := s.serv.UpdateSubscription(id, subs)
	if err != nil {
		return &pb.Subscription_Response{Message: "Error Updating subscription"}, status.Error(codes.Internal, fmt.Sprintf("%s", err))
	}
	return &pb.Subscription_Response{Message: "Success Updating"}, nil

}

func (s *SubService) Delete(ctx context.Context, req *pb.Subscription_Request) (*pb.Subscription_Response, error) {
	id, _ := uuid.Parse(req.Id)
	err := s.serv.DeleteSubscription(id)
	if err != nil {
		return &pb.Subscription_Response{Message: "Error Deleting subscription"}, status.Error(codes.Internal, fmt.Sprintf("%s", err))
	}
	return &pb.Subscription_Response{Message: "success deleting"}, nil

}
