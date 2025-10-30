package endpoint

import (
	"github.com/gin-gonic/gin"
	"http"
	"fmt"
	"github.com/lavatee/camp_backend/internal/model"
)

type SignUpInput struct {
	Email string `json:"email"`
	Password string `json:"password"`
	Name string `json:"name"`
	Tag string `json:"tag"`
	About string `json:"about"`
}

func (e *Endpoint) SignUp(c *gin.Context) {
	var input SignUpInput
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := model.User{
		Email: input.Email,
		PasswordHash: input.Password,
		Name: input.Name,
		Tag: input.Tag,
		About: input.About,
	}
	id, err := e.services.Users.SignUp(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

type SignInInput struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func (e *Endpoint) SignIn(c *gin.Context) {
	var input SignInInput
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	accessToken, refreshToken, err := e.services.Users.SignIn(input.Email, input.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"refresh_token": refreshToken,
	})
}

func (e *Endpoint) CheckTagUnique(c *gin.Context) {
	tag := c.Param("tag")
	if !tag {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("tag is empty")})
		return
	}
	isUnique := e.services.Users.CheckTagUnique(tag)
	c.JSON(http.StatusOK, gin.H{
		"is_unique": isUnique,
	})
}