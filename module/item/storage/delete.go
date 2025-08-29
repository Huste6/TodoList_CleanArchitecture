package storage

import (
	"context"
	"g09/module/item/model"
)

func (s *sqlStore) DeleteItem(ctx context.Context, cond map[string]interface{}) error {
	deletedStatus := "Deleted"
	if err := s.db.Table(model.TodoItem{}.TableName()).Where(cond).Updates(map[string]interface{}{
		"status": deletedStatus,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (s *sqlStore) DeleteItems(ctx context.Context, ids []int) error {
	deletedStatus := "Deleted"
	if err := s.db.Table(model.TodoItem{}.TableName()).Where("id IN ?", ids).Updates(map[string]interface{}{"status": deletedStatus}).Error; err != nil {
		return err
	}
	return nil
}
