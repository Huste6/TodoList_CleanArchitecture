package ginitem

import (
	"g09/common"
	"g09/module/item/biz"
	"g09/module/item/model"
	"g09/module/item/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllItem(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var queryString struct {
			common.Paging
			model.Filter
		}

		if err := c.ShouldBind(&queryString); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		queryString.Paging.Process()

		var res []model.TodoItem

		store := storage.NewSQLStore(db)
		business := biz.NewListItemBiz(store)

		res, err := business.ListItem(c.Request.Context(), &queryString.Filter, &queryString.Paging)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, common.NewSuccessResponse(res, queryString.Paging, queryString.Filter))
	}
}
