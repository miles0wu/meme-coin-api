package service

import (
	"context"
	"github.com/miles0wu/meme-coin-api/internal/domain"
	"github.com/miles0wu/meme-coin-api/internal/repository"
)

var (
	ErrDuplicateName = repository.ErrDuplicateName
	ErrNotFound      = repository.ErrNotFound
)

//go:generate mockgen -source=./coin.go -package=svcmocks -destination=./mocks/coin.mock.go CoinService
type CoinService interface {
	Create(ctx context.Context, coin domain.Coin) (domain.Coin, error)
	Update(ctx context.Context, coin domain.Coin) error
	GetById(ctx context.Context, id int64) (domain.Coin, error)
	DeleteById(ctx context.Context, id int64) error
	IncrPopularityScore(ctx context.Context, id int64) error
}

func NewCoinService(repo repository.CoinRepository) CoinService {
	return &coinService{
		repo: repo,
	}
}

type coinService struct {
	repo repository.CoinRepository
}

func (svc *coinService) Create(ctx context.Context, coin domain.Coin) (domain.Coin, error) {
	return svc.repo.Create(ctx, coin)
}

func (svc *coinService) Update(ctx context.Context, coin domain.Coin) error {
	return svc.repo.Update(ctx, coin)
}

func (svc *coinService) GetById(ctx context.Context, id int64) (domain.Coin, error) {
	return svc.repo.FindById(ctx, id)
}

func (svc *coinService) DeleteById(ctx context.Context, id int64) error {
	return svc.repo.DeleteById(ctx, id)
}

func (svc *coinService) IncrPopularityScore(ctx context.Context, id int64) error {
	return svc.repo.IncrPopularityScore(ctx, id)
}
