package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/miles0wu/meme-coin-api/pkg/logger"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"regexp"
	"testing"
	"time"
)

func TestGormCoinDAO_Insert(t *testing.T) {
	testCases := []struct {
		name    string
		sqlmock func(t *testing.T) *sql.DB

		ctx  context.Context
		coin Coin

		wantErr error
	}{
		{
			name: "insert success",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mockRes := sqlmock.NewResult(1, 1)
				mock.ExpectExec("INSERT INTO `coins` .*").
					WillReturnResult(mockRes)
				return db
			},
			ctx: context.Background(),
			coin: Coin{
				Name: "test",
				Description: sql.NullString{
					String: "test description",
					Valid:  true,
				},
			},
		},
		{
			name: "insert failed - duplicate name",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectExec("INSERT INTO `coins` .*").
					WillReturnError(&mysqlDriver.MySQLError{Number: 1062})
				return db
			},
			ctx:     context.Background(),
			wantErr: ErrDuplicateName,
		},
		{
			name: "insert failed",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectExec("INSERT INTO `coins` .*").
					WillReturnError(errors.New("mock db error"))
				return db
			},
			ctx:     context.Background(),
			wantErr: errors.New("mock db error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB := tc.sqlmock(t)

			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlDB,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true,
			})
			assert.NoError(t, err)
			dao := NewGormCoinDAO(db, logger.NewNopLogger())
			ret, err := dao.Insert(tc.ctx, tc.coin)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.True(t, ret.Id > 0)
			assert.Equal(t, tc.coin.Name, ret.Name)
			assert.Equal(t, tc.coin.Description, ret.Description)
			assert.True(t, ret.CreatedAt > 0)
			assert.True(t, ret.UpdatedAt > 0)
			assert.Equal(t, uint32(0), ret.PopularityScore)
		})
	}
}

func TestGormCoinDAO_UpdateById(t *testing.T) {
	testCases := []struct {
		name    string
		sqlmock func(t *testing.T) *sql.DB

		ctx  context.Context
		coin Coin

		wantErr error
	}{
		{
			name: "update success",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mockRes := sqlmock.NewResult(0, 1)
				mock.ExpectExec("UPDATE `coins` .*").
					WillReturnResult(mockRes)
				return db
			},
			ctx: context.Background(),
			coin: Coin{
				Id:   1,
				Name: "test",
				Description: sql.NullString{
					String: "test description",
					Valid:  true,
				},
			},
		},
		{
			name: "update but no affected",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mockRes := sqlmock.NewResult(0, 0)
				mock.ExpectExec("UPDATE `coins` .*").
					WillReturnResult(mockRes)
				return db
			},
			ctx: context.Background(),
			coin: Coin{
				Id:   1,
				Name: "test",
				Description: sql.NullString{
					String: "test description",
					Valid:  true,
				},
			},
		},
		{
			name: "update failed",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectExec("UPDATE `coins` .*").
					WillReturnError(errors.New("mock db error"))
				return db
			},
			ctx: context.Background(),
			coin: Coin{
				Id:   1,
				Name: "test",
				Description: sql.NullString{
					String: "test description",
					Valid:  true,
				},
			},
			wantErr: errors.New("mock db error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB := tc.sqlmock(t)

			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlDB,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true,
			})
			assert.NoError(t, err)
			dao := NewGormCoinDAO(db, logger.NewNopLogger())
			err = dao.UpdateById(tc.ctx, tc.coin)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestGormCoinDAO_FindById(t *testing.T) {
	nowMs := time.Now().UnixMilli()
	testCases := []struct {
		name    string
		sqlmock func(t *testing.T) *sql.DB

		id int64

		wantRet Coin
		wantErr error
	}{
		{
			name: "success",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `coins` WHERE id = ? ORDER BY `coins`.`id` LIMIT ?")).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at", "popularity_score"}).
						AddRow(1, "test", "test description", nowMs, nowMs, 0))
				return db
			},
			wantRet: Coin{
				Id:   1,
				Name: "test",
				Description: sql.NullString{
					String: "test description",
					Valid:  true,
				},
				CreatedAt:       nowMs,
				UpdatedAt:       nowMs,
				PopularityScore: 0,
			},
			id: 1,
		},
		{
			name: "id not found",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `coins` WHERE id = ? ORDER BY `coins`.`id` LIMIT ?")).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at", "popularity_score"}))
				return db
			},
			id:      1,
			wantErr: ErrRecordNotFound,
		},
		{
			name: "query failed",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `coins` WHERE id = ? ORDER BY `coins`.`id` LIMIT ?")).
					WithArgs(1, 1).
					WillReturnError(errors.New("mock db error"))
				return db
			},
			id:      1,
			wantErr: errors.New("mock db error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB := tc.sqlmock(t)

			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlDB,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true,
			})
			assert.NoError(t, err)
			dao := NewGormCoinDAO(db, logger.NewNopLogger())
			ret, err := dao.FindById(context.Background(), tc.id)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRet, ret)
		})
	}
}

func TestGormCoinDAO_DeleteById(t *testing.T) {
	testCases := []struct {
		name    string
		sqlmock func(t *testing.T) *sql.DB

		id int64

		wantErr error
	}{
		{
			name: "delete success",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mockRes := sqlmock.NewResult(0, 1)
				mock.ExpectExec("DELETE FROM `coins` .*").
					WillReturnResult(mockRes)
				return db
			},
			id: 1,
		},
		{
			name: "delete but no affected",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mockRes := sqlmock.NewResult(0, 0)
				mock.ExpectExec("DELETE FROM `coins` .*").
					WillReturnResult(mockRes)
				return db
			},
			id: 1,
		},
		{
			name: "delete failed",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectExec("DELETE FROM `coins` .*").
					WillReturnError(errors.New("mock db error"))
				return db
			},
			id:      1,
			wantErr: errors.New("mock db error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB := tc.sqlmock(t)

			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlDB,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true,
			})
			assert.NoError(t, err)
			dao := NewGormCoinDAO(db, logger.NewNopLogger())
			err = dao.DeleteById(context.Background(), tc.id)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestGormCoinDAO_IncrPopularityScore(t *testing.T) {
	testCases := []struct {
		name    string
		sqlmock func(t *testing.T) *sql.DB

		id int64

		wantErr error
	}{
		{
			name: "incr success",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mockRes := sqlmock.NewResult(0, 1)
				mock.ExpectExec("UPDATE `coins` .*").
					WillReturnResult(mockRes)
				return db
			},
			id: 1,
		},
		{
			name: "incr but no affected",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mockRes := sqlmock.NewResult(0, 0)
				mock.ExpectExec("UPDATE `coins` .*").
					WillReturnResult(mockRes)
				return db
			},
			id:      1,
			wantErr: ErrRecordNotFound,
		},
		{
			name: "incr failed",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectExec("UPDATE `coins` .*").
					WillReturnError(errors.New("mock db error"))
				return db
			},
			id:      1,
			wantErr: errors.New("mock db error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB := tc.sqlmock(t)

			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlDB,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true,
			})
			assert.NoError(t, err)
			dao := NewGormCoinDAO(db, logger.NewNopLogger())
			err = dao.IncrPopularityScore(context.Background(), tc.id)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
