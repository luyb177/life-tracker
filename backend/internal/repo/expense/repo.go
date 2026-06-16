package expense

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	// Category
	CreateCategory(ctx context.Context, c *Category, tx ...*gorm.DB) error
	FindCategoriesByUser(ctx context.Context, userID uint64, tx ...*gorm.DB) ([]*Category, error)
	FindCategoryByID(ctx context.Context, id uint64, tx ...*gorm.DB) (*Category, error)
	DeleteCategory(ctx context.Context, id uint64, tx ...*gorm.DB) error

	// Log
	CreateLog(ctx context.Context, l *Log, tx ...*gorm.DB) error
	DeleteLog(ctx context.Context, id uint64, tx ...*gorm.DB) error
	FindLogByID(ctx context.Context, id uint64, tx ...*gorm.DB) (*Log, error)
	ListLogsByUser(ctx context.Context, userID uint64, cursorID uint64, cursorTime time.Time, limit int, tx ...*gorm.DB) ([]*Log, error)
	// SumByDate 汇总某用户指定日期的总支出
	SumByDate(ctx context.Context, userID uint64, date time.Time, tx ...*gorm.DB) (float64, error)
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

// ─── Category ───

func (r *repo) CreateCategory(ctx context.Context, c *Category, tx ...*gorm.DB) error {
	return r.getDB(ctx, tx...).Create(c).Error
}

func (r *repo) FindCategoriesByUser(ctx context.Context, userID uint64, tx ...*gorm.DB) ([]*Category, error) {
	var list []*Category
	err := r.getDB(ctx, tx...).Where("user_id = ?", userID).Order("id ASC").Find(&list).Error
	return list, err
}

func (r *repo) FindCategoryByID(ctx context.Context, id uint64, tx ...*gorm.DB) (*Category, error) {
	var c Category
	err := r.getDB(ctx, tx...).Where("id = ?", id).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &c, err
}

func (r *repo) DeleteCategory(ctx context.Context, id uint64, tx ...*gorm.DB) error {
	return r.getDB(ctx, tx...).Delete(&Category{}, id).Error
}

// ─── Log ───

func (r *repo) CreateLog(ctx context.Context, l *Log, tx ...*gorm.DB) error {
	return r.getDB(ctx, tx...).Create(l).Error
}

func (r *repo) DeleteLog(ctx context.Context, id uint64, tx ...*gorm.DB) error {
	return r.getDB(ctx, tx...).Delete(&Log{}, id).Error
}

func (r *repo) FindLogByID(ctx context.Context, id uint64, tx ...*gorm.DB) (*Log, error) {
	var l Log
	err := r.getDB(ctx, tx...).Where("id = ?", id).First(&l).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &l, err
}

// ListLogsByUser 游标分页，按 occurred_at 倒序
func (r *repo) ListLogsByUser(ctx context.Context, userID uint64, cursorID uint64, cursorTime time.Time, limit int, tx ...*gorm.DB) ([]*Log, error) {
	db := r.getDB(ctx, tx...).Model(&Log{}).Where("user_id = ?", userID)

	if cursorID == 0 {
		db = db.Order("occurred_at DESC, id DESC").Limit(limit)
	} else {
		db = db.Where(
			"occurred_at < ? OR (occurred_at = ? AND id < ?)", cursorTime, cursorTime, cursorID,
		).Order("occurred_at DESC, id DESC").Limit(limit)
	}

	var list []*Log
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repo) SumByDate(ctx context.Context, userID uint64, date time.Time, tx ...*gorm.DB) (float64, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var result struct {
		Total float64
	}
	err := r.getDB(ctx, tx...).Model(&Log{}).
		Select("COALESCE(SUM(amount), 0) as total").
		Where("user_id = ? AND occurred_at >= ? AND occurred_at < ?", userID, startOfDay, endOfDay).
		Scan(&result).Error
	return result.Total, err
}
