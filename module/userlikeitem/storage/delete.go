package storage

import (
	"context"
	"g09/common"
	"g09/module/userlikeitem/model"
)

func (s *sqlStore) Delete(ctx context.Context, userId, ItemId int) error {
	var data model.Like

	if err := s.db.Table(data.TableName()).Where("user_id = ? and item_id = ?", userId, ItemId).First(&data).Delete(nil).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
