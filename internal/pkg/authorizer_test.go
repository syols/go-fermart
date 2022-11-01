package pkg

import (
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

const Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk0NjY5MjA2MSwiaXNz" +
	"IjoidXNlcm5hbWUiLCJVc2VybmFtZSI6InVzZXJuYW1lIn0.7Atr6d4ZpCmkGdqrE6yiBfnVktp7wrMURy4WUmRvdXI"

func TestAuthorizerHandler(t *testing.T) {
	auth := Authorizer{
		sign: []byte("some_sign"),
	}
	monkey.Patch(time.Now, func() time.Time {
		return time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
	})
	token, err := auth.CreateToken("username")
	assert.NoError(t, err)
	assert.Equal(t, token, Token)
}

func TestVerifyTokenHandler(t *testing.T) {
	auth := Authorizer{
		sign: []byte("some_sign"),
	}
	createToken, err := auth.CreateToken("username")
	assert.NoError(t, err)

	username, err := auth.VerifyToken(createToken)
	assert.NoError(t, err)
	assert.Equal(t, username, "username")
}
