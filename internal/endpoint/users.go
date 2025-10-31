package endpoint

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lavatee/camp_backend/internal/model"
)

func (e *Endpoint) GetOneUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := e.services.Users.GetOneUser(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

type EditUserDataInput struct {
	Name  string `json:"name"`
	Tag   string `json:"tag"`
	About string `json:"about"`
}

func (e *Endpoint) EditUserData(c *gin.Context) {
	var input EditUserDataInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	newUserInfo := model.User{
		Name:  input.Name,
		Tag:   input.Tag,
		About: input.About,
		Id:    userId,
	}
	if err := e.services.Users.EditUserInfo(newUserInfo); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (e *Endpoint) FindUserByTag(c *gin.Context) {
	tag := c.Param("tag")
	if tag == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("tag is empty")})
		return
	}
	user, err := e.services.Users.FindUserByTag(tag)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (e *Endpoint) NewProfilePhoto(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	url, err := e.services.Users.NewProfilePhoto(userId, file)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"photo_url": url,
	})
}
