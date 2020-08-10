//+build wireinject

package main

import (
	"github.com/dimuska139/golang-api-skeleton/api"
	"github.com/dimuska139/golang-api-skeleton/config"
	"github.com/dimuska139/golang-api-skeleton/database"
	"github.com/dimuska139/golang-api-skeleton/services"
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
)

func InitializeConfig(configPath string) (*config.Config, error) {
	wire.Build(config.NewConfig)
	return &config.Config{}, nil
}

func InitializeDatabase(cfg *config.Config) (*sqlx.DB, error) {
	wire.Build(database.NewDatabase)
	return &sqlx.DB{}, nil
}

func InitializeUsersAPI(db *sqlx.DB) (*api.UsersAPI, error) {
	wire.Build(services.NewUsersService, api.NewUsersAPI)
	return &api.UsersAPI{}, nil
}

func InitializeAuthAPI(cfg *config.Config, db *sqlx.DB) (*api.AuthAPI, error) {
	wire.Build(services.NewUsersService, services.NewAuthService, api.NewAuthAPI)
	return &api.AuthAPI{}, nil
}
