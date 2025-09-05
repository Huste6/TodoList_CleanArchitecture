package storage

import (
	"context"
	"g09/common"
	"g09/module/userlikeitem/model"
	"time"

	"github.com/btcsuite/btcutil/base58"
)

const timeLayout = "2006-01-02T15:04:05.999999"

func (s *sqlStore) ListUsers(ctx context.Context, itemId int, paging *common.Paging) ([]common.SimpleUser, error) {
	var res []model.Like

	db := s.db.Table(model.Like{}.TableName()).Select("user_id").Where("item_id = ?", itemId)
	if err := db.Select("item_id").Count(&paging.Total).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	if v := paging.FakeCursor; v != "" {
		timeCreated, err := time.Parse(timeLayout, string(base58.Decode(v)))
		if err != nil {
			return nil, common.ErrDB(err)
		}
		db = db.Where("created_at < ?", timeCreated.Format("2006-01-02 15:04:05.999999"))
	} else {
		db = db.Offset((paging.Page - 1) * paging.Limit)
	}

	if err := db.
		Select("*").
		Order("created_at desc").
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

	if len(users) > 0 {
		users[len(res)-1].Mask()
		paging.NextCursor = base58.Encode([]byte(users[len(res)-1].CreatedAt.Format(timeLayout)))
	}
	return users, nil
}

func (s *sqlStore) GetItemLikes(ctx context.Context, ids []int) (map[int]int, error) {
	res := make(map[int]int)

	type sqlData struct {
		ItemId int `gorm:"column:item_id"`
		Count  int `gorm:"column:count"`
	}

	var listLike []sqlData
	if err := s.db.Table(model.Like{}.TableName()).Select("item_id, COUNT(item_id) as count").Where("item_id in (?)", ids).Group("item_id").Find(&listLike).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	for _, item := range listLike {
		res[item.ItemId] = item.Count
	}

	return res, nil
}
