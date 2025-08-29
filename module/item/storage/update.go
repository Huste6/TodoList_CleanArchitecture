package storage

import (
	"context"
	"g09/common"
	"g09/module/item/model"
)

func (s *sqlStore) UpdateItem(ctx context.Context, cond map[string]interface{}, dataUpdate *model.TodoItemUpdate) error {
	if err := s.db.Where(cond).Updates(dataUpdate).Error; err != nil {
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
