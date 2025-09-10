package cmd

import (
	"fmt"
	"g09/common"
	"g09/component/rpccaller"
	"g09/component/tokenprovider/jwt"
	"g09/middleware"
	ginitem "g09/module/item/transport/gin"
	"g09/module/upload/transport/ginupload"
	userstorage "g09/module/user/storage"
	ginuser "g09/module/user/transport/gin"
	ginuserlikeitem "g09/module/userlikeitem/transport/gin"
	"g09/pubsub"
	"g09/subscriber"
	"net/http"
	"os"

	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/200Lab-Education/go-sdk/plugin/storage/sdkgorm"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func newService() goservice.Service {
	service := goservice.New(
		goservice.WithName("social-todo-list"),
		goservice.WithVersion("1.0.0"),
		goservice.WithInitRunnable(sdkgorm.NewGormDB("main", common.PluginDBMain)),
		goservice.WithInitRunnable(pubsub.NewPubSub(common.PluginPubSub)),
		goservice.WithInitRunnable(rpccaller.NewApiItemCaller(common.PluginApiItem)),
	)
	return service
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start social TODO service",
	Run: func(cmd *cobra.Command, args []string) {
		systemSecret := os.Getenv("SECRET")
		service := newService()

		serviceLogger := service.Logger("service")
		if err := service.Init(); err != nil {
			serviceLogger.Fatalln(err)
		}

		service.HTTPServer().AddHandler(func(engine *gin.Engine) {
			engine.Use(middleware.Recover())

			db := service.MustGet(common.PluginDBMain).(*gorm.DB)

			authStore := userstorage.NewSQLStore(db)
			tokenProvider := jwt.NewTokenJWTProvider(systemSecret, "jwt")
			middlewareAuth := middleware.RequiredAuth(authStore, tokenProvider)

			v1 := engine.Group("/v1")
			{
				v1.POST("/register", ginuser.Register(db))
				v1.POST("/login", ginuser.Login(db, tokenProvider))
				v1.GET("/profile", middlewareAuth, ginuser.Profile())
				items := v1.Group("/items", middlewareAuth)
				{
					items.POST("", ginitem.CreateItem(db))
					items.POST("/:id/upload", ginupload.UploadAndAttachToItem(db))
					items.GET("", ginitem.GetAllItem(service))
					items.GET("/:id", ginitem.GetItem(db))
					items.PATCH("/:id", ginitem.UpdateItem(db))
					items.PATCH("/status", ginitem.UpdateItems(db))
					items.DELETE("/:id", ginitem.DeleteItem(db))
					items.DELETE("", ginitem.DeleteItems(db))
					items.DELETE("/:id/image", ginupload.DeleteItemImage(db))

					items.POST("/:id/like", ginuserlikeitem.LikeItem(service))
					items.DELETE("/:id/unlike", ginuserlikeitem.UnlikeItem(service))
					items.GET("/:id/liked_users", ginuserlikeitem.ListItem(service))
				}

				rpc := v1.Group("rpc")
				{
					rpc.POST("/get_item_likes", ginuserlikeitem.GetItemLikes(service))
				}

				items = v1.Group("/upload")
				{
					items.POST("", ginupload.Upload(db))
				}
			}

			engine.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "pong",
				})
			})
		})

		_ = subscriber.NewEngine(service).Start()

		if err := service.Start(); err != nil {
			serviceLogger.Fatalln(err)
		}
	},
}

func Execute() {
	rootCmd.AddCommand(outEnvCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
