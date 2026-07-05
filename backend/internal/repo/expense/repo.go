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
	UpdateLog(ctx context.Context, id uint64, updates map[string]interface{}, tx ...*gorm.DB) error
	DeleteLog(ctx context.Context, id uint64, tx ...*gorm.DB) error
	FindLogByID(ctx context.Context, id uint64, tx ...*gorm.DB) (*Log, error)
	RefundLog(ctx context.Context, id uint64, lastUpdatedBy uint64, tx ...*gorm.DB) error
	ListLogsByUser(ctx context.Context, userID uint64, cursorID uint64, cursorTime time.Time, limit int, tx ...*gorm.DB) ([]*Log, error)
	// SumByDate 汇总某用户指定日期的总支出（单位：分）
	SumByDate(ctx context.Context, userID uint64, date time.Time, tx ...*gorm.DB) (int64, error)
	// SumByDateRange 汇总某用户指定日期范围的总支出（单位：分）
	SumByDateRange(ctx context.Context, userID uint64, start, end time.Time, tx ...*gorm.DB) (int64, error)
	// SumByDateRangeGrouped 按分类汇总支出
	SumByDateRangeGrouped(ctx context.Context, userID uint64, start, end time.Time, tx ...*gorm.DB) ([]CategoryTotal, error)
	// ListLogsByDateRange 查询指定日期范围内的支出记录
	ListLogsByDateRange(ctx context.Context, userID uint64, start, end time.Time, tx ...*gorm.DB) ([]*Log, error)
	// CountLogsByCategory 统计某分类下的支出记录数
	CountLogsByCategory(ctx context.Context, categoryID uint64, tx ...*gorm.DB) (int64, error)
	// SumByDay 按天聚合支出
	SumByDay(ctx context.Context, userID uint64, start, end time.Time, tx ...*gorm.DB) ([]DayTotal, error)
	// SumByMonth 按月聚合支出
	SumByMonth(ctx context.Context, userID uint64, start, end time.Time, tx ...*gorm.DB) ([]MonthTotal, error)
}

// CategoryTotal 分类汇总
type CategoryTotal struct {
	CategoryID   uint64 `json:"category_id"`
	CategoryName string `json:"category_name"`
	Total        int64  `json:"total"` // 单位：分
}

// DayTotal 按天汇总
type DayTotal struct {
	Date  string `json:"date"`
	Total int64  `json:"total"` // 单位：分
}

// MonthTotal 按月汇总
type MonthTotal struct {
	Month string `json:"month"` // "2026-01"
	Total int64  `json:"total"` // 单位：分
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
	err := r.getDB(ctx, tx...).Where("user_id = ? OR user_id = 0", userID).Order("id ASC").Find(&list).Error
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

func (r *repo) UpdateLog(ctx context.Context, id uint64, updates map[string]interface{}, tx ...*gorm.DB) error {
	return r.getDB(ctx, tx...).Model(&Log{}).Where("id = ?", id).Updates(updates).Error
}

func (r *repo) DeleteLog(ctx context.Context, id uint64, tx ...*gorm.DB) error {
	return r.getDB(ctx, tx...).Delete(&Log{}, id).Error
}

func (r *repo) RefundLog(ctx context.Context, id uint64, lastUpdatedBy uint64, tx ...*gorm.DB) error {
	now := time.Now()
	return r.getDB(ctx, tx...).Model(&Log{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":          1,
		"refunded_at":     now,
		"last_updated_by": lastUpdatedBy,
		"last_updated_at": now,
	}).Error
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

func (r *repo) SumByDate(ctx context.Context, userID uint64, date time.Time, tx ...*gorm.DB) (int64, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var result struct {
		Total int64
	}
	err := r.getDB(ctx, tx...).Model(&Log{}).
		Select("COALESCE(SUM(amount), 0) as total").
		Where("user_id = ? AND occurred_at >= ? AND occurred_at < ? AND status = 0", userID, startOfDay, endOfDay).
		Scan(&result).Error
	return result.Total, err
}

func (r *repo) SumByDateRange(ctx context.Context, userID uint64, start, end time.Time, tx ...*gorm.DB) (int64, error) {
	var result struct {
		Total int64
	}
	err := r.getDB(ctx, tx...).Model(&Log{}).
		Select("COALESCE(SUM(amount), 0) as total").
		Where("user_id = ? AND occurred_at >= ? AND occurred_at < ? AND status = 0", userID, start, end).
		Scan(&result).Error
	return result.Total, err
}

func (r *repo) ListLogsByDateRange(ctx context.Context, userID uint64, start, end time.Time, tx ...*gorm.DB) ([]*Log, error) {
	var list []*Log
	err := r.getDB(ctx, tx...).
		Where("user_id = ? AND occurred_at >= ? AND occurred_at < ?", userID, start, end).
		Order("occurred_at ASC").
		Find(&list).Error
	return list, err
}

func (r *repo) SumByDateRangeGrouped(ctx context.Context, userID uint64, start, end time.Time, tx ...*gorm.DB) ([]CategoryTotal, error) {
	var results []CategoryTotal
	err := r.getDB(ctx, tx...).Model(&Log{}).
		Select("expense_logs.category_id, expense_categories.name as category_name, COALESCE(SUM(expense_logs.amount), 0) as total").
		Joins("LEFT JOIN expense_categories ON expense_categories.id = expense_logs.category_id AND expense_categories.deleted_at = 0").
		Where("expense_logs.user_id = ? AND expense_logs.occurred_at >= ? AND expense_logs.occurred_at < ? AND expense_logs.status = 0", userID, start, end).
		Group("expense_logs.category_id, expense_categories.name").
		Order("total DESC").
		Scan(&results).Error
	return results, err
}

func (r *repo) SumByMonth(ctx context.Context, userID uint64, start, end time.Time, tx ...*gorm.DB) ([]MonthTotal, error) {
	var results []MonthTotal
	err := r.getDB(ctx, tx...).Model(&Log{}).
		Select("DATE_FORMAT(occurred_at, '%Y-%m') as month, COALESCE(SUM(amount), 0) as total").
		Where("user_id = ? AND occurred_at >= ? AND occurred_at < ? AND status = 0", userID, start, end).
		Group("DATE_FORMAT(occurred_at, '%Y-%m')").
		Order("month ASC").
		Scan(&results).Error
	return results, err
}

func (r *repo) CountLogsByCategory(ctx context.Context, categoryID uint64, tx ...*gorm.DB) (int64, error) {
	var count int64
	err := r.getDB(ctx, tx...).Model(&Log{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count, err
}

func (r *repo) SumByDay(ctx context.Context, userID uint64, start, end time.Time, tx ...*gorm.DB) ([]DayTotal, error) {
	var results []DayTotal
	err := r.getDB(ctx, tx...).Model(&Log{}).
		Select("DATE_FORMAT(occurred_at, '%Y-%m-%d') as date, COALESCE(SUM(amount), 0) as total").
		Where("user_id = ? AND occurred_at >= ? AND occurred_at < ? AND status = 0", userID, start, end).
		Group("DATE_FORMAT(occurred_at, '%Y-%m-%d')").
		Order("date ASC").
		Scan(&results).Error
	return results, err
}
