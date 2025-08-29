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

func CreateItem(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var itemData model.TodoItemCreation
		if err := c.ShouldBind(&itemData); err != nil {
			if appErr, ok := err.(*common.AppError); ok {
				c.JSON(appErr.StatusCode, appErr)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}
		store := storage.NewSQLStore(db)
		business := biz.NewCreateItemBiz(store)
		if err := business.CreateNewItem(c.Request.Context(), &itemData); err != nil {
			if appErr, ok := err.(*common.AppError); ok {
				c.JSON(appErr.StatusCode, appErr)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}
		c.JSON(http.StatusCreated, common.SimpleSuccessResponse(itemData.Id))
	}
}
