package handlers

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/syols/go-devops/internal/models"
	"github.com/syols/go-devops/internal/pkg"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var buf bytes.Buffer
		body, err := ioutil.ReadAll(io.TeeReader(c.Request.Body, &buf))
		if err != nil {
			log.Print(err.Error())
		}
		c.Request.Body = ioutil.NopCloser(&buf)
		log.Print(c.Request.URL.String())
		log.Print(string(body))
		log.Print(headerConvertString(c.Request.Header))
		c.Next()
	}
}

func headerConvertString(h http.Header) string {
	b := new(bytes.Buffer)
	for key, value := range h {
		_, err := fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
		if err != nil {
			return err.Error()
		}
	}
	return b.String()
}

func AuthMiddleware(connection pkg.Database, authorizer pkg.Authorizer) gin.HandlerFunc {
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
			return
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
		context.Set("id", dbUser.ID)
		context.Next()
	}
}
