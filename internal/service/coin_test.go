package service

import (
	"context"
	"errors"
	"github.com/miles0wu/meme-coin-api/internal/domain"
	"github.com/miles0wu/meme-coin-api/internal/repository"
	repomocks "github.com/miles0wu/meme-coin-api/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func Test_coinService_Create(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name string
		mock func(*gomock.Controller) repository.CoinRepository

		coin domain.Coin

		wantRet domain.Coin
		wantErr error
	}{
		{
			name: "create success",
			mock: func(ctrl *gomock.Controller) repository.CoinRepository {
				coinRepo := repomocks.NewMockCoinRepository(ctrl)
				coinRepo.EXPECT().Create(gomock.Any(), domain.Coin{
					Name:        "test",
					Description: "test description",
				}).Return(domain.Coin{
					Id:              1,
					Name:            "test",
					Description:     "test description",
					CreatedAt:       now,
					UpdatedAt:       now,
					PopularityScore: 0,
				}, nil)
				return coinRepo
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
			mock: func(ctrl *gomock.Controller) repository.CoinRepository {
				coinRepo := repomocks.NewMockCoinRepository(ctrl)
				coinRepo.EXPECT().Create(gomock.Any(), domain.Coin{
					Name:        "test",
					Description: "test description",
				}).Return(domain.Coin{}, repository.ErrDuplicateName)
				return coinRepo
			},
			coin: domain.Coin{
				Name:        "test",
				Description: "test description",
			},
			wantRet: domain.Coin{},
			wantErr: repository.ErrDuplicateName,
		},
		{
			name: "db error",
			mock: func(ctrl *gomock.Controller) repository.CoinRepository {
				coinRepo := repomocks.NewMockCoinRepository(ctrl)
				coinRepo.EXPECT().Create(gomock.Any(), domain.Coin{
					Name:        "test",
					Description: "test description",
				}).Return(domain.Coin{}, errors.New("mock db error"))
				return coinRepo
			},
			coin: domain.Coin{
				Name:        "test",
				Description: "test description",
			},
			wantRet: domain.Coin{},
			wantErr: errors.New("mock db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			coinRepo := tc.mock(ctrl)
			svc := NewCoinService(coinRepo)
			ret, err := svc.Create(context.Background(), tc.coin)
			assert.Equal(t, tc.wantErr, err)
			if err == nil {
				return
			}
			assert.Equal(t, tc.wantRet, ret)
		})
	}
}

func Test_coinService_Update(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name string
		mock func(*gomock.Controller) repository.CoinRepository

		coin    domain.Coin
		wantErr error
	}{
		{
			name: "update success",
			mock: func(ctrl *gomock.Controller) repository.CoinRepository {
				coinRepo := repomocks.NewMockCoinRepository(ctrl)
				coinRepo.EXPECT().Update(gomock.Any(), domain.Coin{
					Id:              1,
					Name:            "test",
					Description:     "new test description",
					CreatedAt:       now,
					UpdatedAt:       now,
					PopularityScore: 0,
				}).Return(nil)
				return coinRepo
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
			mock: func(ctrl *gomock.Controller) repository.CoinRepository {
				coinRepo := repomocks.NewMockCoinRepository(ctrl)
				coinRepo.EXPECT().Update(gomock.Any(), domain.Coin{
					Id:              1,
					Name:            "test",
					Description:     "new test description",
					CreatedAt:       now,
					UpdatedAt:       now,
					PopularityScore: 0,
				}).Return(errors.New("mock db error"))
				return coinRepo
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
			coinRepo := tc.mock(ctrl)
			svc := NewCoinService(coinRepo)
			err := svc.Update(context.Background(), tc.coin)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func Test_coinService_GetById(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name string
		mock func(*gomock.Controller) repository.CoinRepository

		id int64

		wantRet domain.Coin
		wantErr error
	}{
		{
			name: "get success",
			mock: func(ctrl *gomock.Controller) repository.CoinRepository {
				coinRepo := repomocks.NewMockCoinRepository(ctrl)
				coinRepo.EXPECT().FindById(gomock.Any(), int64(1)).Return(domain.Coin{
					Id:              1,
					Name:            "test",
					Description:     "test description",
					CreatedAt:       now,
					UpdatedAt:       now,
					PopularityScore: 0,
				}, nil)
				return coinRepo
			},
			id: 1,
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
			name: "id not found",
			mock: func(ctrl *gomock.Controller) repository.CoinRepository {
				coinRepo := repomocks.NewMockCoinRepository(ctrl)
				coinRepo.EXPECT().FindById(gomock.Any(), int64(1)).Return(domain.Coin{}, repository.ErrNotFound)
				return coinRepo
			},
			id:      1,
			wantErr: repository.ErrNotFound,
		},
		{
			name: "db error",
			mock: func(ctrl *gomock.Controller) repository.CoinRepository {
				coinRepo := repomocks.NewMockCoinRepository(ctrl)
				coinRepo.EXPECT().FindById(gomock.Any(), int64(1)).Return(domain.Coin{}, errors.New("mock db error"))
				return coinRepo
			},
			id:      1,
			wantRet: domain.Coin{},
			wantErr: errors.New("mock db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			coinRepo := tc.mock(ctrl)
			svc := NewCoinService(coinRepo)
			ret, err := svc.GetById(context.Background(), tc.id)
			assert.Equal(t, tc.wantErr, err)
			if err == nil {
				return
			}
			assert.Equal(t, tc.wantRet, ret)
		})
	}
}

func Test_coinService_DeleteById(t *testing.T) {
	testCases := []struct {
		name string
		mock func(*gomock.Controller) repository.CoinRepository

		id int64

		wantErr error
	}{
		{
			name: "incr success",
			mock: func(ctrl *gomock.Controller) repository.CoinRepository {
				coinRepo := repomocks.NewMockCoinRepository(ctrl)
				coinRepo.EXPECT().IncrPopularityScore(gomock.Any(), int64(1)).Return(nil)
				return coinRepo
			},
			id:      1,
			wantErr: nil,
		},
		{
			name: "db error",
			mock: func(ctrl *gomock.Controller) repository.CoinRepository {
				coinRepo := repomocks.NewMockCoinRepository(ctrl)
				coinRepo.EXPECT().IncrPopularityScore(gomock.Any(), int64(1)).Return(errors.New("mock db error"))
				return coinRepo
			},
			id:      1,
			wantErr: errors.New("mock db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			coinRepo := tc.mock(ctrl)
			svc := NewCoinService(coinRepo)
			err := svc.IncrPopularityScore(context.Background(), tc.id)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func Test_coinService_IncrPopularityScore(t *testing.T) {
	testCases := []struct {
		name string
		mock func(*gomock.Controller) repository.CoinRepository

		id int64

		wantErr error
	}{
		{
			name: "delete success",
			mock: func(ctrl *gomock.Controller) repository.CoinRepository {
				coinRepo := repomocks.NewMockCoinRepository(ctrl)
				coinRepo.EXPECT().DeleteById(gomock.Any(), int64(1)).Return(nil)
				return coinRepo
			},
			id:      1,
			wantErr: nil,
		},
		{
			name: "db error",
			mock: func(ctrl *gomock.Controller) repository.CoinRepository {
				coinRepo := repomocks.NewMockCoinRepository(ctrl)
				coinRepo.EXPECT().DeleteById(gomock.Any(), int64(1)).Return(errors.New("mock db error"))
				return coinRepo
			},
			id:      1,
			wantErr: errors.New("mock db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			coinRepo := tc.mock(ctrl)
			svc := NewCoinService(coinRepo)
			err := svc.DeleteById(context.Background(), tc.id)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
