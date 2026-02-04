package service

import (
	"errors"
	"subscription/internal/adapters/model"
	"subscription/internal/adapters/repo"

	"github.com/google/uuid"
)

type SubServ interface {
	CreateSubcription(subs model.Subscription) error
	GetSubscriptionByID(id uuid.UUID) ([]model.Subscription, error)
	UpdateSubscription(id uuid.UUID, subs model.Subscription) error
	DeleteSubscription(id uuid.UUID) error
}

type SubService struct {
	repo repo.SubscriptionRepo
}

func NewSubService(repo repo.SubscriptionRepo) SubServ {
	return &SubService{repo: repo}
}

func (s *SubService) CreateSubcription(subs model.Subscription) error {
	if subs.Id_country == uuid.Nil && subs.Id_user == uuid.Nil {
		return errors.New("need id country and user id")
	}

	err := s.repo.CreateSubcription(subs)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubService) GetSubscriptionByID(id uuid.UUID) ([]model.Subscription, error) {
	if id == uuid.Nil {
		return []model.Subscription{}, errors.New("need user id")
	}
	data, err := s.repo.GetBySubscriptionUserID(id)
	if err != nil {
		return []model.Subscription{}, errors.New("no matching data found")
	}
	return data, nil

}

func (s *SubService) UpdateSubscription(id uuid.UUID, subs model.Subscription) error {
	if id == uuid.Nil || subs.Id_country == uuid.Nil {
		return errors.New("needs subscription id and id country")
	}
	err := s.repo.UpdateSubscription(id, subs)
	if err != nil {
		return errors.New("Failed to update subscription")
	}
	return nil
}

func (s *SubService) DeleteSubscription(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("needs subscription id")
	}
	err := s.repo.DeleteSubscription(id)
	if err != nil {
		return errors.New("Failed to delete subscription")
	}
	return nil
}
