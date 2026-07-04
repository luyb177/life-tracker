package lifelog

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, l *LifeLog, tx ...*gorm.DB) error
	Update(ctx context.Context, id uint64, updates map[string]interface{}, tx ...*gorm.DB) error
	Delete(ctx context.Context, id uint64, tx ...*gorm.DB) error
	FindByID(ctx context.Context, id uint64, tx ...*gorm.DB) (*LifeLog, error)
	// ListByUser 游标分页，按 occurred_at 倒序
	ListByUser(ctx context.Context, userID uint64, cursorID uint64, cursorTime time.Time, limit int, tx ...*gorm.DB) ([]*LifeLog, error)
	// ListByUserAndIDs 在指定 ID 范围内游标分页（用于标签过滤后回查）
	ListByUserAndIDs(ctx context.Context, userID uint64, ids []uint64, cursorID uint64, cursorTime time.Time, limit int, tx ...*gorm.DB) ([]*LifeLog, error)
	// FindByDateRange 查询指定日期范围内的生活记录（按 occurred_at 升序）
	FindByDateRange(ctx context.Context, userID uint64, start, end time.Time, tx ...*gorm.DB) ([]*LifeLog, error)
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

func (r *repo) Create(ctx context.Context, l *LifeLog, tx ...*gorm.DB) error {
	return r.getDB(ctx, tx...).Create(l).Error
}

func (r *repo) Update(ctx context.Context, id uint64, updates map[string]interface{}, tx ...*gorm.DB) error {
	return r.getDB(ctx, tx...).Model(&LifeLog{}).Where("id = ?", id).Updates(updates).Error
}

func (r *repo) Delete(ctx context.Context, id uint64, tx ...*gorm.DB) error {
	return r.getDB(ctx, tx...).Delete(&LifeLog{}, id).Error
}

func (r *repo) FindByID(ctx context.Context, id uint64, tx ...*gorm.DB) (*LifeLog, error) {
	var l LifeLog
	err := r.getDB(ctx, tx...).Where("id = ?", id).First(&l).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &l, err
}

// ListByUser 游标分页，按 occurred_at DESC, id DESC
func (r *repo) ListByUser(ctx context.Context, userID uint64, cursorID uint64, cursorTime time.Time, limit int, tx ...*gorm.DB) ([]*LifeLog, error) {
	db := r.getDB(ctx, tx...).Model(&LifeLog{}).Where("user_id = ?", userID)

	if cursorID == 0 {
		db = db.Order("occurred_at DESC, id DESC").Limit(limit)
	} else {
		db = db.Where(
			"occurred_at < ? OR (occurred_at = ? AND id < ?)", cursorTime, cursorTime, cursorID,
		).Order("occurred_at DESC, id DESC").Limit(limit)
	}

	var list []*LifeLog
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// ListByUserAndIDs 在指定 ID 范围内游标分页
func (r *repo) ListByUserAndIDs(ctx context.Context, userID uint64, ids []uint64, cursorID uint64, cursorTime time.Time, limit int, tx ...*gorm.DB) ([]*LifeLog, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	db := r.getDB(ctx, tx...).Model(&LifeLog{}).Where("user_id = ? AND id IN ?", userID, ids)

	if cursorID == 0 {
		db = db.Order("occurred_at DESC, id DESC").Limit(limit)
	} else {
		db = db.Where(
			"occurred_at < ? OR (occurred_at = ? AND id < ?)", cursorTime, cursorTime, cursorID,
		).Order("occurred_at DESC, id DESC").Limit(limit)
	}

	var list []*LifeLog
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// FindByDateRange 查询指定日期范围内的生活记录（按 occurred_at 升序）
func (r *repo) FindByDateRange(ctx context.Context, userID uint64, start, end time.Time, tx ...*gorm.DB) ([]*LifeLog, error) {
	var list []*LifeLog
	err := r.getDB(ctx, tx...).
		Where("user_id = ? AND occurred_at >= ? AND occurred_at < ?", userID, start, end).
		Order("occurred_at ASC").
		Find(&list).Error
	return list, err
}
