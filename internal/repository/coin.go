package repository

import (
	"context"
	"database/sql"
	"github.com/miles0wu/meme-coin-api/internal/domain"
	"github.com/miles0wu/meme-coin-api/internal/repository/cache"
	"github.com/miles0wu/meme-coin-api/internal/repository/dao"
	"github.com/miles0wu/meme-coin-api/pkg/logger"
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
	l     logger.Logger
}

func NewCachedCoinRepository(dao dao.CoinDAO, cache cache.CoinCache, l logger.Logger) CoinRepository {
	return &CachedCoinRepository{
		dao:   dao,
		cache: cache,
		l:     l,
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
	err := repo.dao.UpdateById(ctx, repo.toEntity(coin))
	if err != nil {
		return err
	}
	go func() {
		newCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		er := repo.cache.Del(newCtx, coin.Id)
		if er != nil {
			repo.l.Error("failed to delete coin cache after update coin",
				logger.Int64("coin_id", coin.Id),
				logger.Error(err))
		}
	}()
	return nil
}

func (repo *CachedCoinRepository) FindById(ctx context.Context, id int64) (domain.Coin, error) {
	// get coin from cache, return domain object if hit
	coin, err := repo.cache.Get(ctx, id)
	if err == nil {
		return coin, nil
	}

	// get coin from db
	entity, err := repo.dao.FindById(ctx, id)
	if err != nil {
		return domain.Coin{}, err
	}

	// set coin cache
	coin = repo.toDomain(entity)
	go func() {
		newCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		er := repo.cache.Set(newCtx, coin)
		if er != nil {
			repo.l.Error("failed to set coin cache after get coin from db",
				logger.Int64("coin_id", coin.Id),
				logger.Error(err))
		}
	}()

	return coin, nil
}

func (repo *CachedCoinRepository) DeleteById(ctx context.Context, id int64) error {
	err := repo.dao.DeleteById(ctx, id)
	if err != nil {
		return err
	}
	go func() {
		newCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		er := repo.cache.Del(newCtx, id)
		if er != nil {
			repo.l.Error("failed to delete coin cache after delete coin",
				logger.Int64("coin_id", id),
				logger.Error(err))
		}
	}()
	return err
}

func (repo *CachedCoinRepository) IncrPopularityScore(ctx context.Context, id int64) error {
	err := repo.dao.IncrPopularityScore(ctx, id)
	if err != nil {
		return err
	}
	go func() {
		newCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		er := repo.cache.Del(newCtx, id)
		if er != nil {
			repo.l.Error("failed to delete coin cache after increase popularity score",
				logger.Int64("coin_id", id),
				logger.Error(err))
		}
	}()
	return nil
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
