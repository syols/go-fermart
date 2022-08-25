package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syols/go-devops/internal/models"
	"github.com/syols/go-devops/internal/pkg"
)

func Register(connection pkg.Database, authorizer pkg.Authorizer) gin.HandlerFunc {
	return func(context *gin.Context) {
		user, err := bindUser(context)
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if err := user.Register(context, connection); err != nil {
			context.AbortWithStatus(http.StatusConflict)
			return
		}

		token, err := authorizer.CreateToken(user.Username)
		if err != nil {
			context.AbortWithStatus(http.StatusConflict)
			return
		}

		context.Header("Authorization", "Bearer "+token)
		context.Status(http.StatusOK)
	}
}

func Login(connection pkg.Database, authorizer pkg.Authorizer) gin.HandlerFunc {
	return func(context *gin.Context) {
		user, err := bindUser(context)
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		dbUser, err := user.Login(context, connection)
		if err != nil || dbUser.Username != user.Username {
			context.AbortWithStatus(http.StatusConflict)
			return
		}

		token, err := authorizer.CreateToken(dbUser.Username)
		if err != nil {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		context.Header("Authorization", "Bearer "+token)
		context.Status(http.StatusOK)
	}
}

func bindUser(context *gin.Context) (*models.User, error) {
	var user models.User
	if err := context.BindJSON(&user); err != nil {
		return nil, err
	}
	if err := pkg.Validate(user); err != nil {
		return nil, err
	}
	return &user, nil
}
