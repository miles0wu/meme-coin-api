package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/miles0wu/meme-coin-api/pkg/logger"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicateName  = errors.New("duplicate name")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

//go:generate mockgen -source=./coin.go -package=daomocks -destination=./mocks/coin.mock.go CoinDAO
type CoinDAO interface {
	Insert(ctx context.Context, c Coin) (Coin, error)
	UpdateById(ctx context.Context, entity Coin) error
	FindById(ctx context.Context, uid int64) (Coin, error)
	DeleteById(ctx context.Context, uid int64) error
	IncrPopularityScore(ctx context.Context, id int64) error
}

type GormCoinDAO struct {
	db *gorm.DB
	l  logger.Logger
}

func NewGormCoinDAO(db *gorm.DB, l logger.Logger) CoinDAO {
	return &GormCoinDAO{
		db: db,
		l:  l,
	}
}

func (dao *GormCoinDAO) Insert(ctx context.Context, c Coin) (Coin, error) {
	now := time.Now().UnixMilli()
	c.CreatedAt = now
	c.UpdatedAt = now
	err := dao.db.WithContext(ctx).Create(&c).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			return Coin{}, ErrDuplicateName
		}
	}

	return c, err
}

func (dao *GormCoinDAO) UpdateById(ctx context.Context, entity Coin) error {
	return dao.db.WithContext(ctx).Model(&entity).Where("id = ?", entity.Id).
		Updates(map[string]any{
			"updated_at":  time.Now().UnixMilli(),
			"description": entity.Description,
		}).Error
}

func (dao *GormCoinDAO) FindById(ctx context.Context, id int64) (Coin, error) {
	var res Coin
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&res).Error
	return res, err
}

func (dao *GormCoinDAO) DeleteById(ctx context.Context, id int64) error {
	return dao.db.WithContext(ctx).Where("id = ?", id).Delete(&Coin{}).Error
}

func (dao *GormCoinDAO) IncrPopularityScore(ctx context.Context, id int64) error {
	res := dao.db.WithContext(ctx).Model(&Coin{}).Where("id = ?", id).
		Updates(map[string]any{
			"popularity_score": gorm.Expr("popularity_score + ?", 1),
			"updated_at":       time.Now().UnixMilli(),
		})
	if res.Error == nil && res.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return res.Error
}

type Coin struct {
	Id              int64          `gorm:"primaryKey,autoIncrement"`
	Name            string         `gorm:"unique,type:varchar(255)"`
	Description     sql.NullString `gorm:"type=varchar(128)"`
	CreatedAt       int64
	UpdatedAt       int64
	PopularityScore uint32 `gorm:"default:0"`
}
