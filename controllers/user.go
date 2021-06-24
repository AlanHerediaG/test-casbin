package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/AlanHerediaG/test-casbin/auth"
	"github.com/AlanHerediaG/test-casbin/models"
	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

var tokenManager = auth.TokenManager{}

func Login(c *gin.Context) {
	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	//find user with username
	user, _ := models.UserRepo.FindByID(2)
	//compare the user from the request, with the one we defined:
	if user.UserName != u.UserName || user.Password != u.Password {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}
	ts, err := tokenManager.CreateToken(user.ID, user.UserName)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	c.JSON(http.StatusOK, tokens)
}

func Logout(c *gin.Context) {
	metadata, _ := tokenManager.ExtractTokenMetadata(c.Request)
	if metadata != nil {
		// Delete
	}

	c.JSON(http.StatusOK, "Successfully logged out")
}

func Refresh(c *gin.Context) {
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	refreshToken := mapToken["refresh_token"]

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Refresh token expired")
		return
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userId, roleOk := claims["user_id"].(string)
		if roleOk == false {
			c.JSON(http.StatusUnprocessableEntity, "unauthorized")
			return
		}

		userID, err := strconv.Atoi(userId)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "userId invalid")
			return
		}
		user, err := models.UserRepo.FindByID(userID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "User's not found")
			return
		}

		ts, createErr := tokenManager.CreateToken(userId, user.UserName)
		if createErr != nil {
			c.JSON(http.StatusForbidden, createErr.Error())
			return
		}

		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		c.JSON(http.StatusCreated, tokens)
	} else {
		c.JSON(http.StatusUnauthorized, "refresh expired")
	}
}
