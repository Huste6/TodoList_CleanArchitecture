package subscriber

import (
	"context"
	"g09/common"
	"g09/module/item/storage"
	"g09/pubsub"

	goservice "github.com/200Lab-Education/go-sdk"
	"gorm.io/gorm"
)

func DecreaseLikeCountAfterUserUnLikeItem(serviceCtx goservice.ServiceContext) subJob {
	return subJob{
		Title: "Decrease like count after user unlike item",
		Hld: func(ctx context.Context, message *pubsub.Message) error {
			db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
			data := message.Data().(HasItemId)

			return storage.NewSQLStore(db).DecreaseLikeCount(ctx, data.GetItemId())
		},
	}
}
