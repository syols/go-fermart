package models

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/syols/go-devops/internal/pkg/database"
)

type User struct {
	Id       int    `json:"id" db:"id"`
	Username string `json:"login" db:"login" validate:"min=1"`
	Password string `json:"password" db:"password" validate:"min=1"`
}

func (user User) Validate() error {
	return validator.New().Struct(user)
}

func (user User) Register(ctx context.Context, connection database.Connection) error {
	_, err := connection.Execute(ctx, database.UserRegisterQuery, user)
	return err
}

func (user User) Login(ctx context.Context, connection database.Connection) (*User, error) {
	rows, err := connection.Execute(ctx, database.UserLoginQuery, user)
	if err != nil {
		return nil, err
	}
	return database.ScanOne[User](*rows)
}

func (user User) Verify(ctx context.Context, connection database.Connection) (*User, error) {
	rows, err := connection.Execute(ctx, database.UserSelectQuery, user)
	if err != nil {
		return nil, err
	}
	return database.ScanOne[User](*rows)
}
