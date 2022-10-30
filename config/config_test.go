	package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestConfig(t *testing.T) {
	var configData = `
address:
  host: 0.0.0.0
  port: 8080
accrual: 0.0.0.0:8088
database: "postgres://postgres:postgres@localhost/postgres?sslmode=disable"
sign: "some_sign"
`
	config := Config{}
	assert.NoError(t, yaml.Unmarshal([]byte(configData), &config))
}
