package tag

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	// FindOrCreate 按名称查找标签，不存在则创建（用于前端传入 id=0 的新标签）
	FindOrCreate(ctx context.Context, name string, tx ...*gorm.DB) (*Tag, error)
	// FindByID 按 ID 查找
	FindByID(ctx context.Context, id uint64, tx ...*gorm.DB) (*Tag, error)
	// FindAll 列出所有标签
	FindAll(ctx context.Context, tx ...*gorm.DB) ([]*Tag, error)
	// BatchLink 批量创建 life_log_tags 关联
	BatchLink(ctx context.Context, lifeLogID uint64, tagIDs []uint64, tx ...*gorm.DB) error
	// FindByLifeLogID 查某条 life_log 的所有标签
	FindByLifeLogID(ctx context.Context, lifeLogID uint64, tx ...*gorm.DB) ([]*Tag, error)
	// BatchFindByLifeLogIDs 批量查多条 life_log 的标签
	BatchFindByLifeLogIDs(ctx context.Context, lifeLogIDs []uint64, tx ...*gorm.DB) (map[uint64][]*Tag, error)
	// DeleteByLifeLogID 删除某条 life_log 的所有标签关联
	DeleteByLifeLogID(ctx context.Context, lifeLogID uint64, tx ...*gorm.DB) error
	// FindLifeLogIDsByTagID 按标签 ID 查关联的 life_log ID 列表
	FindLifeLogIDsByTagID(ctx context.Context, tagID uint64, userID uint64, tx ...*gorm.DB) ([]uint64, error)
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

func (r *repo) FindOrCreate(ctx context.Context, name string, tx ...*gorm.DB) (*Tag, error) {
	var t Tag
	err := r.getDB(ctx, tx...).Where("name = ?", name).First(&t).Error
	if err == nil {
		return &t, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	t = Tag{Name: name}
	if err := r.getDB(ctx, tx...).Create(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *repo) FindByID(ctx context.Context, id uint64, tx ...*gorm.DB) (*Tag, error) {
	var t Tag
	err := r.getDB(ctx, tx...).Where("id = ?", id).First(&t).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &t, err
}

func (r *repo) FindAll(ctx context.Context, tx ...*gorm.DB) ([]*Tag, error) {
	var list []*Tag
	err := r.getDB(ctx, tx...).Order("id ASC").Find(&list).Error
	return list, err
}

func (r *repo) BatchLink(ctx context.Context, lifeLogID uint64, tagIDs []uint64, tx ...*gorm.DB) error {
	if len(tagIDs) == 0 {
		return nil
	}
	links := make([]LifeLogTag, 0, len(tagIDs))
	for _, tid := range tagIDs {
		links = append(links, LifeLogTag{LifeLogID: lifeLogID, TagID: tid})
	}
	return r.getDB(ctx, tx...).Create(&links).Error
}

func (r *repo) FindByLifeLogID(ctx context.Context, lifeLogID uint64, tx ...*gorm.DB) ([]*Tag, error) {
	var tags []*Tag
	err := r.getDB(ctx, tx...).
		Joins("JOIN life_log_tags ON life_log_tags.tag_id = tags.id").
		Where("life_log_tags.life_log_id = ?", lifeLogID).
		Find(&tags).Error
	return tags, err
}

func (r *repo) BatchFindByLifeLogIDs(ctx context.Context, lifeLogIDs []uint64, tx ...*gorm.DB) (map[uint64][]*Tag, error) {
	if len(lifeLogIDs) == 0 {
		return nil, nil
	}
	type row struct {
		LifeLogID uint64
		TagID     uint64
		TagName   string
	}
	var rows []row
	err := r.getDB(ctx, tx...).
		Table("life_log_tags").
		Select("life_log_tags.life_log_id, tags.id as tag_id, tags.name as tag_name").
		Joins("JOIN tags ON tags.id = life_log_tags.tag_id").
		Where("life_log_tags.life_log_id IN ?", lifeLogIDs).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make(map[uint64][]*Tag)
	for _, row := range rows {
		result[row.LifeLogID] = append(result[row.LifeLogID], &Tag{ID: row.TagID, Name: row.TagName})
	}
	return result, nil
}

func (r *repo) DeleteByLifeLogID(ctx context.Context, lifeLogID uint64, tx ...*gorm.DB) error {
	return r.getDB(ctx, tx...).Where("life_log_id = ?", lifeLogID).Delete(&LifeLogTag{}).Error
}

func (r *repo) FindLifeLogIDsByTagID(ctx context.Context, tagID uint64, userID uint64, tx ...*gorm.DB) ([]uint64, error) {
	var ids []uint64
	err := r.getDB(ctx, tx...).
		Table("life_log_tags").
		Select("life_log_tags.life_log_id").
		Joins("JOIN life_logs ON life_logs.id = life_log_tags.life_log_id").
		Where("life_log_tags.tag_id = ? AND life_logs.user_id = ?", tagID, userID).
		Pluck("life_log_tags.life_log_id", &ids).Error
	return ids, err
}
