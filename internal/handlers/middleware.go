package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/syols/go-devops/internal/models"
	"github.com/syols/go-devops/internal/pkg/authorizer"
	"github.com/syols/go-devops/internal/pkg/database"
)

func AuthMiddleware(connection database.Connection, authorizer authorizer.Authorizer) gin.HandlerFunc {
	return func(context *gin.Context) {
		auth := context.GetHeader("Authorization")
		if auth == "" {
			context.String(http.StatusUnauthorized, "Authorization header required")
			context.Abort()
			return
		}

		header := strings.Split(auth, " ")
		if len(header) != 2 || header[0] != "Bearer" {
			context.AbortWithStatus(http.StatusUnauthorized)
		}

		token := header[1]
		username, err := authorizer.VerifyToken(token)
		if err != nil {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user := models.User{Username: username}
		dbUser, err := user.Verify(context, connection)
		if err != nil || user.Username != dbUser.Username {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		context.Set("username", user.Username)
		context.Set("id", dbUser.Id)
		context.Next()
	}
}
