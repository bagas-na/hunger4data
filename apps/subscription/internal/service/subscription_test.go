package service

import (
	"errors"
	"subscription/internal/adapters/model"
	mockery "subscription/internal/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateSubscription_Service(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mockery.MockSubscriptionRepo)
		service := NewSubService(mockRepo)
		subs := model.Subscription{
			UserID:      uuid.New(),
			CountryCode: "AFG",
		}

		mockRepo.On("CreateSubcription", subs).Return(nil).Once()

		err := service.CreateSubcription(subs)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Fail - Missing UserID", func(t *testing.T) {
		mockRepo := new(mockery.MockSubscriptionRepo)
		service := NewSubService(mockRepo)
		subs := model.Subscription{
			UserID:      uuid.Nil,
			CountryCode: "USA",
		}

		err := service.CreateSubcription(subs)

		assert.Error(t, err)
		assert.Equal(t, "Must have user_id in jwt", err.Error())
		// Ensure the repo was NOT called
		mockRepo.AssertNotCalled(t, "CreateSubcription", mock.Anything)
	})

	t.Run("Fail - Invalid Country Code", func(t *testing.T) {
		mockRepo := new(mockery.MockSubscriptionRepo)
		service := NewSubService(mockRepo)
		subs := model.Subscription{
			UserID:      uuid.New(),
			CountryCode: "US",
		}

		err := service.CreateSubcription(subs)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "country_code (3 letters)")
		mockRepo.AssertNotCalled(t, "CreateSubcription", mock.Anything)
	})
}

func TestGetSubscriptionByUserID_Flat(t *testing.T) {

	userID := uuid.New()

	t.Run("Success_ValidUser", func(t *testing.T) {
		mockRepo := new(mockery.MockSubscriptionRepo)
		service := &SubService{repo: mockRepo}
		mockData := []model.Subscription{{Id: uuid.New(), UserID: userID}}

		mockRepo.On("GetSubscriptionsByUserID", userID).Return(mockData, nil).Once()

		result, err := service.GetSubscriptionByUserID(userID)

		assert.NoError(t, err)
		assert.Equal(t, mockData, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_NilUserID", func(t *testing.T) {
		mockRepo := new(mockery.MockSubscriptionRepo)
		service := &SubService{repo: mockRepo}
		result, err := service.GetSubscriptionByUserID(uuid.Nil)

		assert.Error(t, err)
		assert.Equal(t, "need user id", err.Error())
		assert.Empty(t, result)
	})

	t.Run("Error_RepoFailure", func(t *testing.T) {
		mockRepo := new(mockery.MockSubscriptionRepo)
		service := &SubService{repo: mockRepo}
		mockRepo.On("GetSubscriptionsByUserID", userID).
			Return([]model.Subscription{}, errors.New("db connection lost")).Once()

		result, err := service.GetSubscriptionByUserID(userID)

		assert.Error(t, err)
		assert.Equal(t, "no matching data found", err.Error())
		assert.Empty(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteSubscription(t *testing.T) {

	userID := uuid.New()
	subsID := uuid.New()

	t.Run("Success_ValidIDs", func(t *testing.T) {
		mockRepo := new(mockery.MockSubscriptionRepo)
		service := NewSubService(mockRepo)
		mockRepo.On("DeleteSubscription", userID, subsID).Return(nil).Once()

		err := service.DeleteSubscription(userID, subsID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_MissingUserID", func(t *testing.T) {
		mockRepo := new(mockery.MockSubscriptionRepo)
		service := NewSubService(mockRepo)
		err := service.DeleteSubscription(uuid.Nil, subsID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing user id")
	})

	t.Run("Error_RepoInternalFailure", func(t *testing.T) {
		mockRepo := new(mockery.MockSubscriptionRepo)
		service := NewSubService(mockRepo)
		repoErr := errors.New("database connection lost")
		mockRepo.On("DeleteSubscription", userID, subsID).Return(repoErr).Once()

		err := service.DeleteSubscription(userID, subsID)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, repoErr))
		assert.Contains(t, err.Error(), "failed to delete subscription")
		mockRepo.AssertExpectations(t)
	})
}
