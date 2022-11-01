package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironment(t *testing.T) {
	address := "0.0.0.0:8080"
	databaseURI := "postgres://uri"
	systemAddress := "0.0.0.0:8088"
	sign := "some_sign"

	t.Setenv("RUN_ADDRESS", address)
	t.Setenv("DATABASE_URI", databaseURI)
	t.Setenv("ACCRUAL_SYSTEM_ADDRESS", systemAddress)
	t.Setenv("SIGN", sign)

	config := Config{}
	config.set(NewEnvironmentVariables().Options()...)
	assert.Equal(t, address, config.ServerAddress.String())
	assert.Equal(t, databaseURI, config.DatabaseURL)
	assert.Equal(t, systemAddress, config.AccrualAddress)
	assert.Equal(t, sign, config.Sign)
}
