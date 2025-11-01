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

func NewEndpoint(services *service.Service) *Endpoint {
	return &Endpoint{
		services: services,
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
		auth.POST("/sign-up", e.SignUp)                //регистрация
		auth.POST("/sign-in", e.SignIn)                // вход в аккаунт
		auth.POST("/refresh", e.Refresh)               //получения новой пары токенов
		auth.GET("/tag-unique/:tag", e.CheckTagUnique) //проверка на уникальность тега (проходит после написания каждой новой буквы в поле ввода тэга на фронтенде)
	}
	api := router.Group("/api", e.Middleware)
	{
		api.GET("/users/:id", e.GetOneUser)           //получение пользователя происходит, когда пользователь нажимает на ник собеседника в чате, либо при заходе в профиль
		api.PUT("/users/:id", e.EditUserData)         //изменение данных в профиле
		api.GET("/users/tag/:tag", e.FindUserByTag)       //при заходе на страницу "/@{tag}" на фронтенде происходит получение пользователя по тегу
		api.POST("/users/photo", e.NewProfilePhoto)   //обновление аватарки
		api.POST("/join-room", e.JoinRoom)            //присоединение к комнате происходит при нажатия кнопки "поиск брата"
		api.POST("/leave-room", e.LeaveRoom)          //выход из поиска брата в чаты пользователя
		api.POST("/next-room", e.NextRoom)            //при нажатии кнопки "пропустить" в поиске брата
		api.GET("/rooms/:id/user", e.GetRoomUser)     //получение данных о пользователе в комнате поиска брата
		api.GET("/user-chats/:query", e.GetUserChats) //получение чатов пользователя и поиск чатов по имени пользователя (query может быть пустой)
		api.GET("/chats/:id/:tz", e.GetOneChat)       //получение одного чата: все сообщения и минимальные данные о собеседнике
		api.PUT("/messages/:id", e.EditMessage)       //изменение сообщения
		api.DELETE("/messages/:id", e.DeleteMessage)  //удаление сообщения
	}
	ws := router.Group("/ws")
	{
		ws.GET("/room/:id/:token", e.RoomWebSocket)     //подключение к комнате поиска брата
		ws.GET("/chat/:id/:token/:tz", e.ChatWebSocket) //подключение к чату
	}
	return router
}
