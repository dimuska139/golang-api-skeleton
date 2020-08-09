//+build wireinject

package main

import (
	"github.com/dimuska139/golang-api-skeleton/api"
	"github.com/dimuska139/golang-api-skeleton/config"
	"github.com/dimuska139/golang-api-skeleton/database"
	"github.com/dimuska139/golang-api-skeleton/services"
	"github.com/google/wire"
)

func InitializeUsersAPI(configPath string) (*api.UsersAPI, error) {
	wire.Build(config.NewConfig, database.NewDatabase, services.NewUsersService, api.NewUsersAPI)
	return &api.UsersAPI{}, nil
}
