package ginitem

import (
	"g09/common"
	"g09/module/item/biz"
	"g09/module/item/model"
	"g09/module/item/repository"
	"g09/module/item/storage"
	usrLikeStore "g09/module/userlikeitem/storage"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
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

		requester := c.MustGet(common.CurrentUser).(common.Requester)

		store := storage.NewSQLStore(db)
		likeStore := usrLikeStore.NewSQLStore(db)
		repo := repository.NewListItemRepo(store, likeStore, requester)
		business := biz.NewListItemBiz(repo, requester)

		res, err := business.ListItem(c.Request.Context(), &queryString.Filter, &queryString.Paging)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		for i := range res {
			res[i].Mask()
		}
		c.JSON(http.StatusOK, common.NewSuccessResponse(res, queryString.Paging, queryString.Filter))
	}
}
