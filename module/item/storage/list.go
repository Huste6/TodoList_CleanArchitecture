package storage

import (
	"context"
	"g09/common"
	"g09/module/item/model"
)

func (s *sqlStore) ListItem(ctx context.Context, filter *model.Filter, paging *common.Paging, moreKeys ...string) ([]model.TodoItem, error) {
	var res []model.TodoItem

	db := s.db.Table(model.TodoItem{}.TableName()).Where("status <> ?", "Deleted")

	// Dùng để lấy ra các todolist của user đấy
	// requester := ctx.Value(common.CurrentUser).(common.Requester)

	// db = db.Where("user_id = ?", requester.GetUserId())

	if f := filter; f != nil {
		if v := f.Status; v != "" {
			db = db.Where("status = ?", v)
		}
	}

	if err := db.Select("id").Count(&paging.Total).Error; err != nil {
		return nil, err
	}

	for i := range moreKeys {
		db = db.Preload(moreKeys[i])
	}

	if err := db.
		Select("*").
		Offset((paging.Page - 1) * paging.Limit).
		Limit(paging.Limit).
		Find(&res).Error; err != nil {
		return nil, common.ErrDB(err)
	}
	return res, nil
}
