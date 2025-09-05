package ginuserlikeitem

import (
	"g09/common"
	itemStore "g09/module/item/storage"
	"g09/module/userlikeitem/biz"
	"g09/module/userlikeitem/model"
	"g09/module/userlikeitem/storage"
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func LikeItem(serviceCtx goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := common.FromBase58(c.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester := c.MustGet(common.CurrentUser).(common.Requester)
		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

		store := storage.NewSQLStore(db)
		itemStore := itemStore.NewSQLStore(db)
		biz := biz.NewUserLikeItemBiz(store, itemStore)

		now := time.Now().UTC()
		if err := biz.LikeItem(c.Request.Context(), &model.Like{
			UserId:    requester.GetUserId(),
			ItemId:    int(id.GetLocalIdlocalId()),
			CreatedAt: &now,
		}); err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
