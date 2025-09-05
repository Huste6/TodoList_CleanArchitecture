package biz

import (
	"context"
	"g09/common"
	"g09/module/userlikeitem/model"
	"log"

	"gorm.io/gorm"
)

type UserUnlikeItemStore interface {
	Find(ctx context.Context, userId, ItemId int) (*model.Like, error)
	Delete(ctx context.Context, userId, ItemId int) error
}
type DecreaseItemStorage interface {
	DecreaseLikeCount(ctx context.Context, id int) error
}

type userUnlikeItemBiz struct {
	store     UserUnlikeItemStore
	itemStore DecreaseItemStorage
}

func NewUserUnlikeItemBiz(store UserUnlikeItemStore, itemStore DecreaseItemStorage) *userUnlikeItemBiz {
	return &userUnlikeItemBiz{store: store, itemStore: itemStore}
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
	go func() {
		defer common.Recover()
		if err := biz.itemStore.DecreaseLikeCount(ctx, ItemId); err != nil {
			log.Println(err.Error())
		}
	}()

	return nil
}
