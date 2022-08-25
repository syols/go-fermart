package models

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/syols/go-devops/internal/pkg/storage"
)

type User struct {
	ID       int    `json:"id" db:"id"`
	Username string `json:"login" db:"login" validate:"min=1"`
	Password string `json:"password" db:"password" validate:"min=1"`
}

func (user User) Validate() error {
	return validator.New().Struct(user)
}

func (user User) Register(ctx context.Context, connection storage.Database) error {
	rows, err := connection.Execute(ctx, "user_register.sql", user)
	if err := rows.Err(); err != nil {
		return err
	}
	return err
}

func (user User) Login(ctx context.Context, connection storage.Database) (*User, error) {
	rows, err := connection.Execute(ctx, "user_login.sql", user)
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return user.bindUser(rows)
}

func (user User) bindUser(rows *sqlx.Rows) (*User, error) {
	var value User
	if rows.Next() {
		if err := rows.StructScan(&value); err != nil {
			return nil, err
		}
	}

	return &value, nil
}

func (user User) Verify(ctx context.Context, connection storage.Database) (*User, error) {
	rows, err := connection.Execute(ctx, "user_select.sql", user)
	if err != nil {
		return nil, err
	}
	return user.bindUser(rows)
}
