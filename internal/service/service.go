package service

import "context"

type Orders interface {
	GetOne(ctx context.Context) (string, error)
	GetMany(ctx context.Context) error
	Remove(ctx context.Context) error
}

type Service struct {
	Repo OrdersRepository
}

func NewService(repo OrdersRepository) Service {
	return Service{repo}
}

func (s Service) GetOne(ctx context.Context) (str string, err error) {

	return s.Repo.GetOne()
}

func (s Service) GetMany(ctx context.Context) (err error) {

	return
}

func (s Service) Remove(ctx context.Context) (err error) {

	return
}
