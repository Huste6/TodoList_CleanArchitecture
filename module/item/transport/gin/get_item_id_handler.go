package ginitem

import (
	"g09/common"
	"g09/module/item/biz"
	"g09/module/item/storage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetItem(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		go func() {
			defer common.Recover()
		}()
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		store := storage.NewSQLStore(db)
		business := biz.NewGetItemBiz(store)

		data, err := business.GetItemById(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		data.Mask()
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data))
	}
}
