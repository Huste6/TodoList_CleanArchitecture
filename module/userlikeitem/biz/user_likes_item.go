package biz

import (
	"context"
	"g09/common"
	"g09/module/userlikeitem/model"
	"g09/pubsub"
	"log"
)

type UserLikeItemStore interface {
	Create(ctx context.Context, data *model.Like) error
}

//	type IncreaseItemStorage interface {
//		IncreaseLikeCount(ctx context.Context, id int) error
//	}
type userLikeItemBiz struct {
	store UserLikeItemStore
	// itemStore IncreaseItemStorage
	ps pubsub.PubSub
}

func NewUserLikeItemBiz(
	store UserLikeItemStore,
	// itemStore IncreaseItemStorage,
	ps pubsub.PubSub,
) *userLikeItemBiz {
	return &userLikeItemBiz{store: store, ps: ps}
}

func (biz *userLikeItemBiz) LikeItem(ctx context.Context, data *model.Like) error {
	if err := biz.store.Create(ctx, data); err != nil {
		return model.ErrCannotLikeItem(err)
	}

	if err := biz.ps.Publish(ctx, common.TopicUserLikedItem, pubsub.NewMessage(data)); err != nil {
		log.Print(err)
	}

	// job := asyncjob.NewJob(func(ctx context.Context) error {
	// 	if err := biz.itemStore.IncreaseLikeCount(ctx, data.ItemId); err != nil {
	// 		return err
	// 	}
	// 	return nil
	// })

	// if err := asyncjob.NewGroup(true, job).Run(ctx); err != nil {
	// 	log.Println(err.Error())

	// }

	return nil
}
