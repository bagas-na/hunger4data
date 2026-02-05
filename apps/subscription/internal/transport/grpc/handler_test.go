package grpcHandler

import (
	"context"
	"errors"
	pb "hunger4data/pb/subcription"
	"subscription/internal/adapters/model"
	mockery "subscription/internal/mocks"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreate_Subscription(t *testing.T) {

	ctx := context.Background()

	t.Run("Success_With_Redis_Interaction", func(t *testing.T) {
		rdb, _ := redismock.NewClientMock()
		mockServ := new(mockery.MockSubServ)
		handler := NewSubHand(mockServ, rdb)
		userID := uuid.New().String()
		req := &pb.Subscription_Request{
			UserId:      userID,
			CountryCode: "CAF",
		}

		mockServ.On("CreateSubcription", mock.Anything).Return(nil).Once()

		resp, err := handler.Create_Subscription(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, "Subscription successful", resp.Message)

	})
	t.Run("Failure - Repository Returns Error", func(t *testing.T) {
		rdb, _ := redismock.NewClientMock()
		mockServ := new(mockery.MockSubServ)
		handler := NewSubHand(mockServ, rdb)
		userID := uuid.New().String()
		req := &pb.Subscription_Request{
			UserId:      userID,
			CountryCode: "AFG",
		}

		mockServ.On("CreateSubcription", mock.Anything).
			Return(errors.New("db error"))

		res, err := handler.Create_Subscription(context.Background(), req)

		assert.Error(t, err)
		assert.Equal(t, "Error Creating subscription", res.Message)

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
	})
}

func TestCreate_Get_Subscription(t *testing.T) {

	ctx := context.Background()

	t.Run("Success - Returns list of subscriptions", func(t *testing.T) {
		rdb, _ := redismock.NewClientMock()
		mockServ := new(mockery.MockSubServ)
		handler := NewSubHand(mockServ, rdb)
		userID := uuid.New()
		req := &pb.Subscription_Request{UserId: userID.String()}

		mockData := []model.Subscription{
			{Id: uuid.New(), UserId: userID, CountryCode: "US"},
			{Id: uuid.New(), UserId: userID, CountryCode: "UK"},
		}

		mockServ.On("GetSubscriptionByUserID", userID).Return(mockData, nil)
		res, err := handler.Get_Subscriptions(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, "Success Getting Subscription", res.Message)
		assert.Len(t, res.Subscription, 2)

		assert.Equal(t, mockData[0].Id.String(), res.Subscription[0].Id)
		assert.Equal(t, mockData[1].CountryCode, res.Subscription[1].CountryCode)

		mockServ.AssertExpectations(t)
	})

	t.Run("Failure - Repository Error", func(t *testing.T) {
		rdb, _ := redismock.NewClientMock()
		mockServ := new(mockery.MockSubServ)
		handler := NewSubHand(mockServ, rdb)
		userID := uuid.New()
		req := &pb.Subscription_Request{UserId: userID.String()}
		mockServ.On("GetSubscriptionByUserID", mock.Anything).
			Return(nil, errors.New("not found"))

		// Act
		res, err := handler.Get_Subscriptions(context.Background(), req)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "Error Getting subscription", res.Message)

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
	})
}

func TestCreate_Delete_Subscription(t *testing.T) {

	ctx := context.Background()

	t.Run("Success - Valid IDs", func(t *testing.T) {
		rdb, _ := redismock.NewClientMock()
		mockServ := new(mockery.MockSubServ)
		handler := NewSubHand(mockServ, rdb)
		userID := uuid.New()
		subID := uuid.New()
		req := &pb.Subscription_Request{
			UserId: userID.String(),
			Id:     subID.String(),
		}

		mockServ.On("DeleteSubscription", userID, subID).Return(nil)

		res, err := handler.Delete_Subscription(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, "success deleting", res.Message)
		mockServ.AssertExpectations(t)
	})

	t.Run("Failure - Invalid Subscription ID Format", func(t *testing.T) {
		rdb, _ := redismock.NewClientMock()
		mockServ := new(mockery.MockSubServ)
		handler := NewSubHand(mockServ, rdb)
		userID := uuid.New()
		req := &pb.Subscription_Request{
			UserId: userID.String(),
			Id:     "not-a-uuid",
		}
		res, err := handler.Delete_Subscription(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, res.Message, "Error Parsing SubcpritionId")

		st, _ := status.FromError(err)
		assert.Equal(t, codes.Internal, st.Code())
		mockServ.AssertNotCalled(t, "DeleteSubscription", mock.Anything, mock.Anything)
	})

	t.Run("Failure - Repository Deletion Error", func(t *testing.T) {
		rdb, _ := redismock.NewClientMock()
		mockServ := new(mockery.MockSubServ)
		handler := NewSubHand(mockServ, rdb)
		userID := uuid.New()
		subID := uuid.New()
		req := &pb.Subscription_Request{UserId: userID.String(), Id: subID.String()}

		mockServ.On("DeleteSubscription", userID, subID).
			Return(errors.New("db error"))

		res, err := handler.Delete_Subscription(context.Background(), req)

		assert.Error(t, err)
		assert.Equal(t, "Error Deleting subscription", res.Message)
		mockServ.AssertExpectations(t)
	})
}
