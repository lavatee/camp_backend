package endpoint

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lavatee/camp_backend/internal/service"
)

type Endpoint struct {
	services *service.Service
}

func NewEndpoint(services service.Service) *Endpoint {
	return &Endpoint{
		services: &services,
	}
}

func (e *Endpoint) InitRoutes() *gin.Engine {
	router := gin.New()
	config := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(config))
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", e.SignUp)
		auth.POST("/sign-in", e.SignIn)
		auth.POST("/refresh", e.Refresh)
		auth.GET("/tag-unique/:tag", e.CheckTagUnique)
	}
	api := router.Group("/api", e.Middleware)
	{
		api.GET("/users/:id", e.GetOneUser)
		api.PUT("/users/:id", e.EditUserData)
		api.GET("/users/:tag", e.FindUserByTag)
		api.POST("/users/photo", e.NewProfilePhoto)
		api.POST("/join-room", e.JoinRoom)
		api.POST("/leave-room", e.LeaveRoom)
		api.POST("/next-room", e.NextRoom)
		api.GET("/rooms/:id/user", e.GetRoomUser)
		api.POST("/chats", e.PostChat)
		api.GET("/user-chats/:query", e.GetUserChats)
		api.GET("/chats/:id", e.GetOneChat)
		api.PUT("/messages/:id", e.EditMessage)
		api.DELETE("/messages/:id", e.DeleteMessage)
	}
	ws := router.Group("/ws")
	{
		
	}
	return router
}
