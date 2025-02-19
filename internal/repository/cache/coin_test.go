package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/miles0wu/meme-coin-api/internal/domain"
	"github.com/miles0wu/meme-coin-api/internal/repository/cache/redismocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestRedisCoinCache_Set(t *testing.T) {
	coin := domain.Coin{
		Id: 1,
	}

	keyFunc := func(id int64) string {
		return fmt.Sprintf("coin:%d", id)
	}
	testCases := []struct {
		name string
		mock func(*gomock.Controller) redis.Cmdable

		ctx  context.Context
		coin domain.Coin

		wantErr error
	}{
		{
			name: "set success",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				bs, err := json.Marshal(coin)
				assert.NoError(t, err)
				cmd := redismocks.NewMockCmdable(ctrl)
				mockRes := redis.NewStatusResult("OK", nil)
				cmd.EXPECT().Set(gomock.Any(), keyFunc(coin.Id), bs, 15*time.Minute).Return(mockRes)
				return cmd
			},
			ctx:  context.Background(),
			coin: coin,
		},
		{
			name: "redis conn error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				bs, err := json.Marshal(coin)
				assert.NoError(t, err)
				cmd := redismocks.NewMockCmdable(ctrl)
				mockRes := redis.NewStatusResult("", errors.New("redis conn error"))
				cmd.EXPECT().Set(gomock.Any(), keyFunc(coin.Id), bs, 15*time.Minute).Return(mockRes)
				return cmd
			},
			ctx:     context.Background(),
			coin:    coin,
			wantErr: errors.New("redis conn error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cmd := tc.mock(ctrl)
			cache := NewRedisCoinCache(cmd)

			err := cache.Set(tc.ctx, tc.coin)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestRedisCoinCache_Get(t *testing.T) {
	coin := domain.Coin{
		Id: 1,
	}

	keyFunc := func(id int64) string {
		return fmt.Sprintf("coin:%d", id)
	}
	testCases := []struct {
		name string
		mock func(*gomock.Controller) redis.Cmdable

		ctx context.Context
		id  int64

		wantErr error
	}{
		{
			name: "get success",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				bs, err := json.Marshal(coin)
				assert.NoError(t, err)
				cmd := redismocks.NewMockCmdable(ctrl)
				mockRes := redis.NewStringResult(string(bs), nil)
				cmd.EXPECT().Get(gomock.Any(), keyFunc(1)).Return(mockRes)
				return cmd
			},
			ctx: context.Background(),
			id:  1,
		},
		{
			name: "key not found",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				mockRes := redis.NewStringResult("", redis.Nil)
				cmd.EXPECT().Get(gomock.Any(), keyFunc(2)).Return(mockRes)
				return cmd
			},
			ctx:     context.Background(),
			id:      2,
			wantErr: redis.Nil,
		},
		{
			name: "redis conn error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				mockRes := redis.NewStringResult("", errors.New("redis conn error"))
				cmd.EXPECT().Get(gomock.Any(), keyFunc(coin.Id)).Return(mockRes)
				return cmd
			},
			ctx:     context.Background(),
			id:      1,
			wantErr: errors.New("redis conn error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cmd := tc.mock(ctrl)
			cache := NewRedisCoinCache(cmd)

			_, err := cache.Get(tc.ctx, tc.id)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestRedisCoinCache_Del(t *testing.T) {
	keyFunc := func(id int64) string {
		return fmt.Sprintf("coin:%d", id)
	}
	testCases := []struct {
		name string
		mock func(*gomock.Controller) redis.Cmdable

		ctx context.Context
		id  int64

		wantErr error
	}{
		{
			name: "delete success",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				mockRes := redis.NewIntResult(1, nil)
				cmd.EXPECT().Del(gomock.Any(), keyFunc(1)).Return(mockRes)
				return cmd
			},
			ctx: context.Background(),
			id:  1,
		},
		{
			name: "key not found",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				mockRes := redis.NewIntResult(0, nil)
				cmd.EXPECT().Del(gomock.Any(), keyFunc(2)).Return(mockRes)
				return cmd
			},
			ctx:     context.Background(),
			id:      2,
			wantErr: nil,
		},
		{
			name: "redis conn error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				mockRes := redis.NewIntResult(0, errors.New("redis conn error"))
				cmd.EXPECT().Del(gomock.Any(), keyFunc(1)).Return(mockRes)
				return cmd
			},
			ctx:     context.Background(),
			id:      1,
			wantErr: errors.New("redis conn error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cmd := tc.mock(ctrl)
			cache := NewRedisCoinCache(cmd)

			err := cache.Del(tc.ctx, tc.id)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
