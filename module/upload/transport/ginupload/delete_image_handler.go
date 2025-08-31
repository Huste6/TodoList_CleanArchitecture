package ginupload

import (
	"g09/common"
	"g09/module/item/biz"
	"g09/module/item/storage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DeleteItemImage(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		itemId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
			return
		}

		// Dùng storage để lấy item
		itemStore := storage.NewSQLStore(db)
		item, err := itemStore.GetItem(c.Request.Context(), map[string]interface{}{"id": itemId})
		if err != nil {
			if appErr, ok := err.(*common.AppError); ok {
				c.JSON(appErr.StatusCode, appErr)
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			}
			return
		}

		// Kiểm tra có image không
		if item.Image == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Item has no image to delete"})
			return
		}

		oldImageURL := item.Image.Url

		// Xóa image từ database bằng business layer
		requester := c.MustGet(common.CurrentUser).(common.Requester)

		itemBiz := biz.NewUpdateItemBiz(itemStore, requester)

		if err := itemBiz.DeleteItemImage(c.Request.Context(), itemId); err != nil {
			if appErr, ok := err.(*common.AppError); ok {
				c.JSON(appErr.StatusCode, appErr)
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image from database"})
			}
			return
		}

		// Xóa từ Cloudinary bằng helper function
		if err := DeleteOldImage(c.Request.Context(), oldImageURL); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"data":    true,
				"warning": "Image deleted from database but failed to delete from Cloudinary",
			})
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
