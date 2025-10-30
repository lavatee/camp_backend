package endpoint

import (
	"github.com/gin-gonic/gin"
	"http"
	"github.com/lavatee/camp_backend/internal/model"
)

type JoinRoomInput struct {
	UserId int `json:"user_id"`
}

func (e *Endpoint) JoinRoom(c *gin.Context) {
	var input JoinRoomInput
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	room, err := e.services.Rooms.JoinRoom(input.UserId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"room": room,
	})
}

type LeaveRoomInput struct {
	UserId int `json:"user_id"`
	RoomId int `json:"room_id"`
}

func (e *Endpoint) LeaveRoom(c *gin.Context) {
	var input LeaveRoomInput
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := e.services.Rooms.LeaveRoom(input.UserId, input.RoomId); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (e *Endpoint) NextRoom(c *gin.Context) {
	var input LeaveRoomInput
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	room, err := e.services.Rooms.NextRoom(input.UserId, input.RoomId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"room": room,
	})
}

func (e *Endpoint) GetRoomUser()