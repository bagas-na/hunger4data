package service

import (
	"errors"
	"fmt"
	"subscription/internal/adapters/model"
	"subscription/internal/adapters/repo"

	"github.com/google/uuid"
)

type SubServ interface {
	CreateSubcription(subs model.Subscription) error
	GetSubscriptionByUserID(userId uuid.UUID) ([]model.Subscription, error)
	// UpdateSubscription(id uuid.UUID, subs model.Subscription) error
	DeleteSubscription(userId uuid.UUID, subsId uuid.UUID) error
}

type SubService struct {
	repo repo.SubscriptionRepo
}

func NewSubService(repo repo.SubscriptionRepo) SubServ {
	return &SubService{repo: repo}
}

func (s *SubService) CreateSubcription(subs model.Subscription) error {
	if subs.UserId == uuid.Nil {
		return errors.New("Must have user_id in jwt")
	}

	if subs.CountryCode == "" || len(subs.CountryCode) != 3 {
		return errors.New("Must have country_code (3 letters)")
	}

	err := s.repo.CreateSubcription(subs)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubService) GetSubscriptionByUserID(userId uuid.UUID) ([]model.Subscription, error) {
	if userId == uuid.Nil {
		return []model.Subscription{}, errors.New("need user id")
	}
	data, err := s.repo.GetSubscriptionsByUserID(userId)
	if err != nil {
		return []model.Subscription{}, errors.New("no matching data found")
	}
	return data, nil

}

// func (s *SubService) UpdateSubscription(id uuid.UUID, subs model.Subscription) error {
// 	if subs.UserId == uuid.Nil {
// 		return errors.New("Must have user_id in jwt")
// 	}

// 	if subs.CountryCode == "" {
// 		return errors.New("Must have country_code (3 letters)")
// 	}

//		err := s.repo.UpdateSubscription(id, subs)
//		if err != nil {
//			return errors.New("Failed to update subscription")
//		}
//		return nil
//	}

func (s *SubService) DeleteSubscription(userId uuid.UUID, subsId uuid.UUID) error {
	if userId == uuid.Nil {
		return errors.New("missing user id")
	}

	if subsId == uuid.Nil {
		return errors.New("missing subscription id")
	}

	err := s.repo.DeleteSubscription(userId, subsId)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return nil
}
