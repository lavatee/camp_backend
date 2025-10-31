package endpoint

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (e *Endpoint) GetUserChats(c *gin.Context) {
	query := c.Param("query")
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	chats, err := e.services.Chats.GetUserChats(userId, query)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"chats": chats,
	})
}

func (e *Endpoint) GetOneChat(c *gin.Context) {
	chatId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	timeZone := c.Param("tz")
	if timeZone == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("time zone is empty")})
		return
	}
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	chat, isTreeLegit, err := e.services.Chats.GetOneChat(userId, chatId, timeZone)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"chat":          chat,
		"is_tree_legit": isTreeLegit,
	})
}

type EditMessageInput struct {
	NewText string `json:"text"`
}

func (e *Endpoint) EditMessage(c *gin.Context) {
	var input EditMessageInput
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	messageId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if err := e.services.Chats.EditMessage(messageId, userId, input.NewText); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (e *Endpoint) DeleteMessage(c *gin.Context) {
	messageId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if err := e.services.Chats.DeleteMessage(messageId, userId); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
