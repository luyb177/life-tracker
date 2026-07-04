package summary

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// Repository summary 仓储接口
type Repository interface {
	Create(ctx context.Context, s *Summary, tx ...*gorm.DB) error
	Update(ctx context.Context, id uint64, updates map[string]interface{}, tx ...*gorm.DB) error
	FindByID(ctx context.Context, id uint64, tx ...*gorm.DB) (*Summary, error)
	FindByPeriod(ctx context.Context, userID uint64, periodType uint8, periodStart string, tx ...*gorm.DB) (*Summary, error)
	// FindByPeriodAndSource 按周期+来源查询（用于 AI 去重）
	FindByPeriodAndSource(ctx context.Context, userID uint64, periodType uint8, periodStart string, source uint8, tx ...*gorm.DB) (*Summary, error)
	// ExistsByPeriodAndSource 检查指定周期+来源是否已有记录
	ExistsByPeriodAndSource(ctx context.Context, userID uint64, periodType uint8, periodStart string, source uint8, tx ...*gorm.DB) (bool, error)
	// FindByPeriodRange 查询某周期类型在时间范围内的总结（用于周报聚合日报、月报聚合周报等）
	FindByPeriodRange(ctx context.Context, userID uint64, periodType uint8, start, end string, tx ...*gorm.DB) ([]*Summary, error)
	// FindByPeriodRangeAndSource 按周期类型+时间范围+来源查询
	FindByPeriodRangeAndSource(ctx context.Context, userID uint64, periodType uint8, start, end string, source uint8, tx ...*gorm.DB) ([]*Summary, error)
	// ListByUser 游标分页，按 created_at 倒序
	ListByUser(ctx context.Context, userID uint64, periodType uint8, cursorID uint64, cursorTime time.Time, limit int, tx ...*gorm.DB) ([]*Summary, error)
	Delete(ctx context.Context, id uint64, tx ...*gorm.DB) error
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db: db}
}

func (r *repo) getDB(ctx context.Context, tx ...*gorm.DB) *gorm.DB {
	if len(tx) > 0 && tx[0] != nil {
		return tx[0]
	}
	return r.db.WithContext(ctx)
}

func (r *repo) Create(ctx context.Context, s *Summary, tx ...*gorm.DB) error {
	return r.getDB(ctx, tx...).Create(s).Error
}

func (r *repo) Update(ctx context.Context, id uint64, updates map[string]interface{}, tx ...*gorm.DB) error {
	return r.getDB(ctx, tx...).Model(&Summary{}).Where("id = ?", id).Updates(updates).Error
}

func (r *repo) FindByID(ctx context.Context, id uint64, tx ...*gorm.DB) (*Summary, error) {
	var s Summary
	err := r.getDB(ctx, tx...).Where("id = ?", id).First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &s, err
}

func (r *repo) FindByPeriod(ctx context.Context, userID uint64, periodType uint8, periodStart string, tx ...*gorm.DB) (*Summary, error) {
	var s Summary
	err := r.getDB(ctx, tx...).
		Where("user_id = ? AND period_type = ? AND period_start = ?", userID, periodType, periodStart).
		Order("created_at DESC").
		First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &s, err
}

func (r *repo) FindByPeriodAndSource(ctx context.Context, userID uint64, periodType uint8, periodStart string, source uint8, tx ...*gorm.DB) (*Summary, error) {
	var s Summary
	err := r.getDB(ctx, tx...).
		Where("user_id = ? AND period_type = ? AND period_start = ? AND source = ?", userID, periodType, periodStart, source).
		First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &s, err
}

func (r *repo) ExistsByPeriodAndSource(ctx context.Context, userID uint64, periodType uint8, periodStart string, source uint8, tx ...*gorm.DB) (bool, error) {
	var count int64
	err := r.getDB(ctx, tx...).Model(&Summary{}).
		Where("user_id = ? AND period_type = ? AND period_start = ? AND source = ?", userID, periodType, periodStart, source).
		Count(&count).Error
	return count > 0, err
}

func (r *repo) FindByPeriodRange(ctx context.Context, userID uint64, periodType uint8, start, end string, tx ...*gorm.DB) ([]*Summary, error) {
	var list []*Summary
	err := r.getDB(ctx, tx...).
		Where("user_id = ? AND period_type = ? AND period_start >= ? AND period_start < ?", userID, periodType, start, end).
		Order("period_start ASC").
		Find(&list).Error
	return list, err
}

func (r *repo) FindByPeriodRangeAndSource(ctx context.Context, userID uint64, periodType uint8, start, end string, source uint8, tx ...*gorm.DB) ([]*Summary, error) {
	var list []*Summary
	err := r.getDB(ctx, tx...).
		Where("user_id = ? AND period_type = ? AND period_start >= ? AND period_start < ? AND source = ?", userID, periodType, start, end, source).
		Order("period_start ASC").
		Find(&list).Error
	return list, err
}

func (r *repo) ListByUser(ctx context.Context, userID uint64, periodType uint8, cursorID uint64, cursorTime time.Time, limit int, tx ...*gorm.DB) ([]*Summary, error) {
	db := r.getDB(ctx, tx...).Model(&Summary{}).Where("user_id = ?", userID)
	if periodType > 0 {
		db = db.Where("period_type = ?", periodType)
	}

	if cursorID == 0 {
		db = db.Order("created_at DESC, id DESC").Limit(limit)
	} else {
		db = db.Where(
			"created_at < ? OR (created_at = ? AND id < ?)", cursorTime, cursorTime, cursorID,
		).Order("created_at DESC, id DESC").Limit(limit)
	}

	var list []*Summary
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repo) Delete(ctx context.Context, id uint64, tx ...*gorm.DB) error {
	return r.getDB(ctx, tx...).Delete(&Summary{}, id).Error
}
