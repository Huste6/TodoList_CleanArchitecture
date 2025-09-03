package storage

import (
	"context"
	"g09/common"
	"g09/module/userlikeitem/model"
)

func (s *sqlStore) ListUsers(ctx context.Context, itemId int, paging *common.Paging) ([]common.SimpleUser, error) {
	var res []model.Like

	db := s.db.Table(model.Like{}.TableName()).Select("user_id").Where("item_id = ?", itemId)
	if err := db.Select("item_id").Count(&paging.Total).Error; err != nil {
		return nil, common.ErrDB(err)
	}
	if err := db.
		Select("*").
		Offset((paging.Page - 1) * paging.Limit).
		Limit(paging.Limit).
		Preload("User").
		Find(&res).Error; err != nil {
		return nil, common.ErrDB(err)
	}
	users := make([]common.SimpleUser, len(res))
	for i := range users {
		users[i] = *res[i].User
		users[i].UpdatedAt = nil
		users[i].CreatedAt = res[i].CreatedAt
	}
	return users, nil
}
