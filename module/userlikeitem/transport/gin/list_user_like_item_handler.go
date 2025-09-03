package ginuserlikeitem

import (
	"g09/common"
	"g09/module/userlikeitem/biz"
	"g09/module/userlikeitem/storage"
	"net/http"

	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListItem(serviceCtx goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := common.FromBase58(c.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var queryString struct {
			common.Paging
		}

		if err := c.ShouldBind(&queryString); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		queryString.Process()
		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

		store := storage.NewSQLStore(db)
		biz := biz.NewListUserLikeItemBiz(store)

		res, err := biz.ListUserLikeItem(c.Request.Context(), int(id.GetLocalIdlocalId()), &queryString.Paging)
		if err != nil {
			panic(err)
		}
		for i := range res {
			res[i].Mask()
		}
		c.JSON(http.StatusOK, common.NewSuccessResponse(res, queryString.Paging, nil))
	}
}
