package subscriber

import (
	"context"
	"g09/common"
	"g09/module/item/storage"
	"g09/pubsub"

	goservice "github.com/200Lab-Education/go-sdk"
	"gorm.io/gorm"
)

type HasItemId interface {
	GetItemId() int
}

// func IcreaseLikeCountAfterUserLikeItem(serviceCtx goservice.ServiceContext, ctx context.Context) {
// 	ps := serviceCtx.MustGet(common.PluginPubSub).(pubsub.PubSub)
// 	db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

// 	c, _ := ps.Subscribe(ctx, common.TopicUserLikedItem)

// 	go func() {
// 		defer common.Recover()
// 		for msg := range c {
// 			data := msg.Data().(*model.Like)
// 			if err := storage.NewSQLStore(db).IncreaseLikeCount(ctx, data.GetItemId()); err != nil {
// 				log.Println(err)
// 			}
// 		}
// 	}()
// }

func IcreaseLikeCountAfterUserLikeItem(serviceCtx goservice.ServiceContext) subJob {
	return subJob{
		Title: "Increase like count after user like item",
		Hld: func(ctx context.Context, message *pubsub.Message) error {
			db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
			data := message.Data().(HasItemId)

			return storage.NewSQLStore(db).IncreaseLikeCount(ctx, data.GetItemId())
		},
	}
}
