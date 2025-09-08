package subscriber

import (
	"context"
	"g09/pubsub"
	"log"

	goservice "github.com/200Lab-Education/go-sdk"
)

type HasUserId interface {
	GetUserId() int
}

func PushNotificationAfterUserLikeItem(serviceCtx goservice.ServiceContext) subJob {
	return subJob{
		Title: "Push notification after user likes item",
		Hld: func(ctx context.Context, message *pubsub.Message) error {
			data := message.Data().(HasUserId)
			log.Println("Push notification to user id: ", data.GetUserId())
			return nil
		},
	}
}
