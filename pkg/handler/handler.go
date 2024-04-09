package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vekshinnikita/golang_music/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"*"}
	router.Use(cors.New(corsConfig))

	auth := router.Group("/auth")
	{
		auth.POST("/signup", h.singUp)
		auth.POST("/signin", h.singIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		track := api.Group("/track")
		{
			track.POST("", h.AddTrack)
			track.PATCH("", h.UpdateTrack)
			track.DELETE("/:id", h.DeleteTrack)
			track.GET("/:id", h.GetTrackInfo)
			track.GET("/:id/streaming", h.StreamingTrack)
			track.GET("/:id/poster", h.GetPoster)
			track.POST("/upload", h.UploadTrack)
			track.POST("/poster/upload", h.UploadPoster)

		}
	}

	return router
}
