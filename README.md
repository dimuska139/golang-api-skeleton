# DEPRECATED! See [go-api-layout](https://github.com/dimuska139/go-api-layout) please

# Golang API skeleton

* Config file reading ([viper](github.com/spf13/viper))
* Migrations ([golang-migrate](https://github.com/golang-migrate/migrate))
* Compile-time Dependency Injection ([google/wire](https://github.com/google/wire))
* Working with database (where is no [GORM](http://gorm.io/index.html) in this skeleton but you can easily integrate it)
* Token-Based Authentication (with sliding sessions)

Tests will be soon :)

## Migrations

1. [Create migration file](https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md) in `/migrations` directory 
(also you can use [CLI](https://github.com/golang-migrate/migrate#cli-usage) for it).
1. Build your application.
1. Migrations applying automatically after you run compiled binary file.

## Dependency injection

[google/wire](https://github.com/google/wire) - DI without magic and run-time reflection.

To append new service to initialization you should:

1. Create service and "provider" for it (see **NewUsersAPI** in `/api/users.go` for example).
1. Inject provider to initialization in `wire.go` (first line `//+build wireinject` is definitely needed)
1. Run `wire` command to generate `wire_gen.go` (file with generated initialization steps)
1. Build/Run your app

Also you can read detailed [tutorial](https://github.com/google/wire/blob/master/_tutorial/README.md).
