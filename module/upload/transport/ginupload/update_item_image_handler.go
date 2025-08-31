package ginupload

import (
	"g09/common"
	"g09/module/item/biz"
	"g09/module/item/model"
	"g09/module/item/storage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UploadAndAttachToItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		itemId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
			return
		}

		// Get current item to check for existing image
		itemStore := storage.NewSQLStore(db)
		currentItem, err := itemStore.GetItem(c.Request.Context(), map[string]interface{}{"id": itemId})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}

		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer file.Close()

		// Process file upload using helper function
		img, err := ProcessFileUpload(c, file, header)
		if err != nil {
			if appErr, ok := err.(*common.AppError); ok {
				c.JSON(appErr.StatusCode, appErr)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		// Delete old image if exists
		if currentItem.Image != nil && currentItem.Image.Url != "" {
			if err := DeleteOldImage(c.Request.Context(), currentItem.Image.Url); err != nil {
				return
			}
		}

		// Gán ảnh mới cho item
		requester := c.MustGet(common.CurrentUser).(common.Requester)

		itemBiz := biz.NewUpdateItemBiz(itemStore, requester)
		updateData := &model.TodoItemUpdate{Image: img}

		if err := itemBiz.UpdateItemById(c.Request.Context(), itemId, updateData); err != nil {
			if appErr, ok := err.(*common.AppError); ok {
				c.JSON(appErr.StatusCode, appErr)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(img))
	}
}
