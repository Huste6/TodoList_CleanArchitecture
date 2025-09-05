package storage

import (
	"context"
	"errors"
	"g09/common"
	"g09/module/item/model"

	"gorm.io/gorm"
)

func (s *sqlStore) UpdateItem(ctx context.Context, cond map[string]interface{}, dataUpdate *model.TodoItemUpdate) error {
	if err := s.db.Where(cond).Updates(dataUpdate).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}

func (s *sqlStore) DeleteItemImage(ctx context.Context, itemId int) error {
	if err := s.db.Model(&model.TodoItem{}).Where("id = ?", itemId).Update("image", nil).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}

func (s *sqlStore) UpdateItems(ctx context.Context, ids []int, status string) error {
	if err := s.db.Table(model.TodoItemUpdate{}.TableName()).Where("ID in ?", ids).Updates(map[string]interface{}{"status": status}).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}

func (s *sqlStore) IncreaseLikeCount(ctx context.Context, id int) error {
	db := s.db.Table(model.TodoItem{}.TableName())

	if err := db.Where("id = ?", id).Update("like_count", gorm.Expr("like_count + ?", 1)).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}

func (s *sqlStore) DecreaseLikeCount(ctx context.Context, id int) error {
	db := s.db.Table(model.TodoItem{}.TableName())

	var item model.TodoItem
	if err := db.Where("id = ?", id).First(&item).Error; err != nil {
		return common.ErrDB(err)
	}
	if item.LikedCount <= 0 {
		return common.ErrInvalidRequest(errors.New("like count must be greater than 0"))
	}
	if err := db.Where("id = ?", id).Update("like_count", gorm.Expr("like_count - ?", 1)).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}
