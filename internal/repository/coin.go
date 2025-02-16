package repository

import (
	"context"
	"database/sql"
	"github.com/miles0wu/meme-coin-api/internal/domain"
	"github.com/miles0wu/meme-coin-api/internal/repository/cache"
	"github.com/miles0wu/meme-coin-api/internal/repository/dao"
	"time"
)

var (
	ErrDuplicateName = dao.ErrDuplicateName
	ErrNotFound      = dao.ErrRecordNotFound
)

//go:generate mockgen -source=./coin.go -package=repomocks -destination=./mocks/coin.mock.go CoinRepository
type CoinRepository interface {
	Create(ctx context.Context, coin domain.Coin) (domain.Coin, error)
	Update(ctx context.Context, coin domain.Coin) error
	FindById(ctx context.Context, id int64) (domain.Coin, error)
	DeleteById(ctx context.Context, id int64) error
	IncrPopularityScore(ctx context.Context, id int64) error
}

type CachedCoinRepository struct {
	dao   dao.CoinDAO
	cache cache.CoinCache
}

// NewCachedCoinRepository
// TODO: add cache mechanism
func NewCachedCoinRepository(dao dao.CoinDAO) CoinRepository {
	return &CachedCoinRepository{
		dao: dao,
	}
}

func (repo *CachedCoinRepository) Create(ctx context.Context, coin domain.Coin) (domain.Coin, error) {
	dc, err := repo.dao.Insert(ctx, repo.toEntity(coin))
	if err != nil {
		return domain.Coin{}, err
	}
	return repo.toDomain(dc), nil
}

func (repo *CachedCoinRepository) Update(ctx context.Context, coin domain.Coin) error {
	return repo.dao.UpdateById(ctx, repo.toEntity(coin))
}

func (repo *CachedCoinRepository) FindById(ctx context.Context, id int64) (domain.Coin, error) {
	coin, err := repo.dao.FindById(ctx, id)
	if err != nil {
		return domain.Coin{}, err
	}

	return repo.toDomain(coin), nil
}

func (repo *CachedCoinRepository) DeleteById(ctx context.Context, id int64) error {
	return repo.dao.DeleteById(ctx, id)
}

func (repo *CachedCoinRepository) IncrPopularityScore(ctx context.Context, id int64) error {
	return repo.dao.IncrPopularityScore(ctx, id)
}

func (repo *CachedCoinRepository) toEntity(c domain.Coin) dao.Coin {
	return dao.Coin{
		Id:          c.Id,
		Name:        c.Name,
		Description: sql.NullString{String: c.Description, Valid: c.Description != ""},
	}
}

func (repo *CachedCoinRepository) toDomain(c dao.Coin) domain.Coin {
	return domain.Coin{
		Id:              c.Id,
		Name:            c.Name,
		Description:     c.Description.String,
		CreatedAt:       time.UnixMilli(c.CreatedAt),
		UpdatedAt:       time.UnixMilli(c.UpdatedAt),
		PopularityScore: c.PopularityScore,
	}
}
