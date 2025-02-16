package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/miles0wu/meme-coin-api/internal/domain"
	"github.com/miles0wu/meme-coin-api/internal/repository/dao"
	daomocks "github.com/miles0wu/meme-coin-api/internal/repository/dao/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCachedCoinRepository_Create(t *testing.T) {
	nowMs := time.Now().UnixMilli()
	now := time.UnixMilli(nowMs)
	testCases := []struct {
		name string
		mock func(*gomock.Controller) dao.CoinDAO

		coin domain.Coin

		wantRet domain.Coin
		wantErr error
	}{
		{
			name: "create success",
			mock: func(ctrl *gomock.Controller) dao.CoinDAO {
				coinDAO := daomocks.NewMockCoinDAO(ctrl)
				coinDAO.EXPECT().Insert(gomock.Any(), dao.Coin{
					Name:        "test",
					Description: sql.NullString{String: "test description", Valid: true},
				}).Return(dao.Coin{
					Id:              1,
					Name:            "test",
					Description:     sql.NullString{String: "test description", Valid: true},
					CreatedAt:       nowMs,
					UpdatedAt:       nowMs,
					PopularityScore: 0,
				}, nil)
				return coinDAO
			},
			coin: domain.Coin{
				Name:        "test",
				Description: "test description",
			},
			wantRet: domain.Coin{
				Id:              1,
				Name:            "test",
				Description:     "test description",
				CreatedAt:       now,
				UpdatedAt:       now,
				PopularityScore: 0,
			},
			wantErr: nil,
		},
		{
			name: "duplicate name error",
			mock: func(ctrl *gomock.Controller) dao.CoinDAO {
				coinDAO := daomocks.NewMockCoinDAO(ctrl)
				coinDAO.EXPECT().Insert(gomock.Any(), dao.Coin{
					Name:        "test",
					Description: sql.NullString{String: "test description", Valid: true},
				}).Return(dao.Coin{}, dao.ErrDuplicateName)
				return coinDAO
			},
			coin: domain.Coin{
				Name:        "test",
				Description: "test description",
			},
			wantErr: dao.ErrDuplicateName,
		},
		{
			name: "db error",
			mock: func(ctrl *gomock.Controller) dao.CoinDAO {
				coinDAO := daomocks.NewMockCoinDAO(ctrl)
				coinDAO.EXPECT().Insert(gomock.Any(), dao.Coin{
					Name:        "test",
					Description: sql.NullString{String: "test description", Valid: true},
				}).Return(dao.Coin{}, errors.New("mock db error"))
				return coinDAO
			},
			coin: domain.Coin{
				Name:        "test",
				Description: "test description",
			},
			wantErr: errors.New("mock db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			coinDAO := tc.mock(ctrl)
			repo := NewCachedCoinRepository(coinDAO)
			ret, err := repo.Create(context.Background(), tc.coin)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRet, ret)
		})
	}
}

func TestCachedCoinRepository_Update(t *testing.T) {
	nowMs := time.Now().UnixMilli()
	now := time.UnixMilli(nowMs)
	testCases := []struct {
		name string
		mock func(*gomock.Controller) dao.CoinDAO

		coin domain.Coin

		wantErr error
	}{
		{
			name: "update success",
			mock: func(ctrl *gomock.Controller) dao.CoinDAO {
				coinDAO := daomocks.NewMockCoinDAO(ctrl)
				coinDAO.EXPECT().UpdateById(gomock.Any(), dao.Coin{
					Id:          1,
					Name:        "test",
					Description: sql.NullString{String: "new test description", Valid: true},
				}).Return(nil)
				return coinDAO
			},
			coin: domain.Coin{
				Id:              1,
				Name:            "test",
				Description:     "new test description",
				CreatedAt:       now,
				UpdatedAt:       now,
				PopularityScore: 0,
			},
			wantErr: nil,
		},
		{
			name: "db error",
			mock: func(ctrl *gomock.Controller) dao.CoinDAO {
				coinDAO := daomocks.NewMockCoinDAO(ctrl)
				coinDAO.EXPECT().UpdateById(gomock.Any(), dao.Coin{
					Id:          1,
					Name:        "test",
					Description: sql.NullString{String: "new test description", Valid: true},
				}).Return(errors.New("mock db error"))
				return coinDAO
			},
			coin: domain.Coin{
				Id:              1,
				Name:            "test",
				Description:     "new test description",
				CreatedAt:       now,
				UpdatedAt:       now,
				PopularityScore: 0,
			},
			wantErr: errors.New("mock db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			coinDAO := tc.mock(ctrl)
			repo := NewCachedCoinRepository(coinDAO)
			err := repo.Update(context.Background(), tc.coin)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestCachedCoinRepository_FindById(t *testing.T) {
	nowMs := time.Now().UnixMilli()
	now := time.UnixMilli(nowMs)
	testCases := []struct {
		name string
		mock func(*gomock.Controller) dao.CoinDAO

		id int64

		wantRet domain.Coin
		wantErr error
	}{
		{
			name: "find success",
			mock: func(ctrl *gomock.Controller) dao.CoinDAO {
				coinDAO := daomocks.NewMockCoinDAO(ctrl)
				coinDAO.EXPECT().FindById(gomock.Any(), int64(1)).Return(dao.Coin{
					Id:              1,
					Name:            "test",
					Description:     sql.NullString{String: "new test description", Valid: true},
					CreatedAt:       nowMs,
					UpdatedAt:       nowMs,
					PopularityScore: 0,
				}, nil)
				return coinDAO
			},
			id: 1,
			wantRet: domain.Coin{
				Id:              1,
				Name:            "test",
				Description:     "new test description",
				CreatedAt:       now,
				UpdatedAt:       now,
				PopularityScore: 0,
			},
			wantErr: nil,
		},
		{
			name: "record not found",
			mock: func(ctrl *gomock.Controller) dao.CoinDAO {
				coinDAO := daomocks.NewMockCoinDAO(ctrl)
				coinDAO.EXPECT().FindById(gomock.Any(), int64(1)).Return(dao.Coin{}, dao.ErrRecordNotFound)
				return coinDAO
			},
			id:      1,
			wantErr: dao.ErrRecordNotFound,
		},
		{
			name: "db error",
			mock: func(ctrl *gomock.Controller) dao.CoinDAO {
				coinDAO := daomocks.NewMockCoinDAO(ctrl)
				coinDAO.EXPECT().FindById(gomock.Any(), int64(1)).Return(dao.Coin{}, errors.New("mock db error"))
				return coinDAO
			},
			id:      1,
			wantErr: errors.New("mock db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			coinDAO := tc.mock(ctrl)
			repo := NewCachedCoinRepository(coinDAO)
			ret, err := repo.FindById(context.Background(), tc.id)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRet, ret)
		})
	}
}

func TestCachedCoinRepository_DeleteById(t *testing.T) {
	testCases := []struct {
		name string
		mock func(*gomock.Controller) dao.CoinDAO

		id int64

		wantErr error
	}{
		{
			name: "delete success",
			mock: func(ctrl *gomock.Controller) dao.CoinDAO {
				coinDAO := daomocks.NewMockCoinDAO(ctrl)
				coinDAO.EXPECT().DeleteById(gomock.Any(), int64(1)).Return(nil)
				return coinDAO
			},
			id:      1,
			wantErr: nil,
		},
		{
			name: "db error",
			mock: func(ctrl *gomock.Controller) dao.CoinDAO {
				coinDAO := daomocks.NewMockCoinDAO(ctrl)
				coinDAO.EXPECT().DeleteById(gomock.Any(), int64(1)).Return(errors.New("mock db error"))
				return coinDAO
			},
			id:      1,
			wantErr: errors.New("mock db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			coinDAO := tc.mock(ctrl)
			repo := NewCachedCoinRepository(coinDAO)
			err := repo.DeleteById(context.Background(), tc.id)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestCachedCoinRepository_IncrPopularityScore(t *testing.T) {
	testCases := []struct {
		name string
		mock func(*gomock.Controller) dao.CoinDAO

		id int64

		wantErr error
	}{
		{
			name: "incr success",
			mock: func(ctrl *gomock.Controller) dao.CoinDAO {
				coinDAO := daomocks.NewMockCoinDAO(ctrl)
				coinDAO.EXPECT().IncrPopularityScore(gomock.Any(), int64(1)).Return(nil)
				return coinDAO
			},
			id:      1,
			wantErr: nil,
		},
		{
			name: "0 row affected",
			mock: func(ctrl *gomock.Controller) dao.CoinDAO {
				coinDAO := daomocks.NewMockCoinDAO(ctrl)
				coinDAO.EXPECT().IncrPopularityScore(gomock.Any(), int64(1)).Return(dao.ErrRecordNotFound)
				return coinDAO
			},
			id:      1,
			wantErr: dao.ErrRecordNotFound,
		},
		{
			name: "db error",
			mock: func(ctrl *gomock.Controller) dao.CoinDAO {
				coinDAO := daomocks.NewMockCoinDAO(ctrl)
				coinDAO.EXPECT().IncrPopularityScore(gomock.Any(), int64(1)).Return(errors.New("mock db error"))
				return coinDAO
			},
			id:      1,
			wantErr: errors.New("mock db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			coinDAO := tc.mock(ctrl)
			repo := NewCachedCoinRepository(coinDAO)
			err := repo.IncrPopularityScore(context.Background(), tc.id)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
