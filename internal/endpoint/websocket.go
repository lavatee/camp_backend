package endpoint

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lavatee/camp_backend/internal/model"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Room struct {
	clients            map[int]*websocket.Conn
	FriendshipRequests map[int]bool
}

type WsMessage struct {
	Text          string `json:"text"`
	UserId        int    `json:"user_id"`
	MessageId     int    `json:"message_id"`
	IsTreeUpdated bool   `json:"is_tree_updated"`
}

var rooms = make(map[int]*Room)
var chats = make(map[int]*Room)

func (e *Endpoint) RoomWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()
	token := c.Param("token")
	claims, err := e.services.Users.ParseToken(token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	floatId, ok := claims["id"].(float64)
	if !ok {
		fmt.Println("invalid type of id")
		return
	}
	userId := int(floatId)
	roomId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if _, ok := rooms[roomId]; !ok {
		rooms[roomId] = &Room{
			clients: map[int]*websocket.Conn{
				userId: conn,
			},
		}
	} else {
		rooms[roomId].clients[userId] = conn
		data, err := json.Marshal(WsMessage{Text: fmt.Sprintf("event:user_connect:%d", userId), UserId: 0}) //сообщение о том, что второй пользователь подключился (id = 0 обозначает, что сообщение пришло от сервера), на фронтеде идет получение второго пользователя комнаты у обоих пользователей и теперь пользователи могут обмениваться сообщениями
		if err != nil {
			delete(rooms[roomId].clients, userId)
			return
		}
		for _, client := range rooms[roomId].clients {
			err = client.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				delete(rooms[roomId].clients, userId)
				break
			}
		}

	}

	conn.SetCloseHandler(func(code int, text string) error {
		delete(rooms[roomId].clients, userId)
		return nil
	})
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			delete(rooms[roomId].clients, userId)
			break
		}
		if string(msg) == "*friendship*" { //сигнал, который приходит, когда пользователь нажимает на кнопку "Добавить в друзья"
			rooms[roomId].FriendshipRequests[userId] = true
			if len(rooms[roomId].FriendshipRequests) == 2 {
				friends := make([]int, 0)
				for id, _ := range rooms[roomId].clients {
					friends = append(friends, id)
				}
				if len(friends) < 2 {
					logrus.Errorf("chat can be created if 2 users in a room")
					return
				}
				chatId, err := e.services.Chats.CreateChat(friends[0], friends[1])
				if err != nil {
					logrus.Errorf("error while creating chat: %s", err.Error())
					continue
				}
				data, err := json.Marshal(WsMessage{Text: fmt.Sprintf("event:new_chat:%d", chatId), UserId: 0}) //сообщение о том, что чат создан (id = 0 обозначает, что сообщение пришло от сервера), на фронтеде появляется кнопка "перейти в чат"
				if err != nil {
					delete(rooms[roomId].clients, userId)
					break
				}
				for _, client := range rooms[roomId].clients {
					err = client.WriteMessage(websocket.TextMessage, data)
					if err != nil {
						delete(rooms[roomId].clients, userId)
						break
					}
				}
				delete(rooms[roomId].clients, userId)
			}
		}
		data, err := json.Marshal(WsMessage{Text: string(msg), UserId: userId})
		if err != nil {
			delete(rooms[roomId].clients, userId)
			break
		}
		for _, client := range rooms[roomId].clients {
			err = client.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				delete(rooms[roomId].clients, userId)
				break
			}
		}
	}
}

func (e *Endpoint) ChatWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()
	timeZone := c.Param("tz")
	token := c.Param("token")
	claims, err := e.services.Users.ParseToken(token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	floatId, ok := claims["id"].(float64)
	if !ok {
		fmt.Println("invalid type of id")
		return
	}
	userId := int(floatId)
	chatId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if _, ok := chats[chatId]; !ok {
		chats[chatId] = &Room{
			clients: map[int]*websocket.Conn{
				userId: conn,
			},
		}
	} else {
		chats[chatId].clients[userId] = conn
	}

	conn.SetCloseHandler(func(code int, text string) error {
		delete(chats[chatId].clients, userId)
		return nil
	})
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			delete(chats[chatId].clients, userId)
			break
		}
		isRead := false
		if len(chats[chatId].clients) == 2 {
			isRead = true
		}
		message := model.Message{
			Text:   string(msg),
			UserId: userId,
			IsRead: isRead,
			ChatId: chatId,
		}
		messageId, isTreeUpdated, err := e.services.Chats.CreateMessage(message, timeZone)
		if err != nil {
			delete(chats[chatId].clients, userId)
			break
		}
		data, err := json.Marshal(WsMessage{Text: string(msg), UserId: userId, MessageId: messageId, IsTreeUpdated: isTreeUpdated})
		if err != nil {
			delete(chats[chatId].clients, userId)
			break
		}
		for _, client := range chats[chatId].clients {
			err = client.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				delete(chats[chatId].clients, userId)
				break
			}
		}
	}
}
