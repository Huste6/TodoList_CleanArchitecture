package biz

import (
	"context"
	"g09/common"
	"g09/module/userlikeitem/model"

	"gorm.io/gorm"
)

type UserUnlikeItemStore interface {
	Find(ctx context.Context, userId, ItemId int) (*model.Like, error)
	Delete(ctx context.Context, userId, ItemId int) error
}

type userUnlikeItemBiz struct {
	store UserUnlikeItemStore
}

func NewUserUnlikeItemBiz(store UserUnlikeItemStore) *userUnlikeItemBiz {
	return &userUnlikeItemBiz{store: store}
}

func (biz *userUnlikeItemBiz) UnLikeItem(ctx context.Context, userId, ItemId int) error {
	_, err := biz.store.Find(ctx, userId, ItemId)
	if err == gorm.ErrRecordNotFound {
		return common.RecordNotFound
	}
	if err != nil {
		return model.ErrCannotUnlikeItem(err)
	}

	if err := biz.store.Delete(ctx, userId, ItemId); err != nil {
		return model.ErrCannotLikeItem(err)
	}
	return nil
}
