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

func (s *SubService) Create_Subscription(ctx context.Context, req *pb.Subscription_Request) (*pb.Subscription_Response, error) {
	userId, _ := uuid.Parse(req.UserId)
	subs := model.Subscription{
		UserId:      userId,
		CountryCode: req.CountryCode,
	}
	err := s.serv.CreateSubcription(subs)
	if err != nil {
		return &pb.Subscription_Response{Message: "Error Creating subscription"}, status.Error(codes.Internal, fmt.Sprintf("%s", err))
	}

	return &pb.Subscription_Response{Message: "Subscription successful"}, nil
}

func (s *SubService) Get_Subscriptions(ctx context.Context, req *pb.Subscription_Request) (*pb.Get_Subscriptions_Response, error) {
	userId, _ := uuid.Parse(req.UserId)
	data, err := s.serv.GetSubscriptionByUserID(userId)
	if err != nil {
		return &pb.Get_Subscriptions_Response{Message: "Error Getting subscription"}, status.Error(codes.Internal, fmt.Sprintf("%s", err))
	}
	protoSubs := []*pb.Subscription{}
	for _, sub := range data {
		protoSubs = append(protoSubs, &pb.Subscription{
			Id:          sub.Id.String(),
			UserId:      sub.UserId.String(),
			CountryCode: sub.CountryCode,
		})
	}
	return &pb.Get_Subscriptions_Response{Subscription: protoSubs, Message: "Success Getting Subscription"}, nil

}

// func (s *SubService) Update_Subscription(ctx context.Context, req *pb.Subscription_Request) (*pb.Subscription_Response, error) {
// 	subs := model.Subscription{
// 		CountryCode: req.CountryCode,
// 	}
// 	id, _ := uuid.Parse(req.Id)
// 	err := s.serv.UpdateSubscription(id, subs)
// 	if err != nil {
// 		return &pb.Subscription_Response{Message: "Error Updating subscription"}, status.Error(codes.Internal, fmt.Sprintf("%s", err))
// 	}
// 	return &pb.Subscription_Response{Message: "Success Updating"}, nil

// }

func (s *SubService) Delete_Subscription(ctx context.Context, req *pb.Subscription_Request) (*pb.Subscription_Response, error) {
	userId, _ := uuid.Parse(req.UserId)
	subscriptionId, err := uuid.Parse(req.Id)
	if err != nil {
		return &pb.Subscription_Response{
				Message: "Error Parsing SubcpritionId. Expect uuid",
			},
			status.Error(codes.Internal, fmt.Sprintf("%s", err))
	}

	err = s.serv.DeleteSubscription(userId, subscriptionId)
	if err != nil {
		return &pb.Subscription_Response{
				Message: "Error Deleting subscription",
			},
			status.Error(codes.Internal, fmt.Sprintf("%s", err))
	}
	return &pb.Subscription_Response{
		Message: "success deleting",
	}, nil

}
