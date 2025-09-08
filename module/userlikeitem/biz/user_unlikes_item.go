package biz

import (
	"context"
	"g09/common"
	"g09/module/userlikeitem/model"
	"g09/pubsub"
	"log"

	"gorm.io/gorm"
)

type UserUnlikeItemStore interface {
	Find(ctx context.Context, userId, ItemId int) (*model.Like, error)
	Delete(ctx context.Context, userId, ItemId int) error
}

// type DecreaseItemStorage interface {
// 	DecreaseLikeCount(ctx context.Context, id int) error
// }

type userUnlikeItemBiz struct {
	store UserUnlikeItemStore
	ps    pubsub.PubSub
}

func NewUserUnlikeItemBiz(store UserUnlikeItemStore, ps pubsub.PubSub) *userUnlikeItemBiz {
	return &userUnlikeItemBiz{store: store, ps: ps}
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

	if err := biz.ps.Publish(ctx, common.TopicUserUnLikedItem, pubsub.NewMessage(&model.Like{UserId: userId, ItemId: ItemId})); err != nil {
		log.Print(err)
	}

	// job := asyncjob.NewJob(func(ctx context.Context) error {
	// 	if err := biz.itemStore.DecreaseLikeCount(ctx, ItemId); err != nil {
	// 		return err
	// 	}
	// 	return nil
	// })

	// if err := asyncjob.NewGroup(true, job).Run(ctx); err != nil {
	// 	log.Println(err.Error())
	// }

	return nil
}
