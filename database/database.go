package database

import (
	"github.com/dimuska139/golang-api-skeleton/config"
	_ "github.com/dimuska139/golang-api-skeleton/migrations"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"log"
	"time"
)

func NewDatabase(config *config.Config) (*sqlx.DB, error) {
	connConfig := pgx.ConnConfig{
		Host:     config.Database.Host,
		Port:     config.Database.Port,
		Database: config.Database.Name,
		User:     config.Database.User,
		Password: config.Database.Password,
	}

	connPool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     connConfig,
		AfterConnect:   nil,
		MaxConnections: config.Database.MaxConnections,
		AcquireTimeout: 30 * time.Second,
	})
	if err != nil {
		return nil, errors.Wrap(err, "connection pool initialization failed")
	}

	nativeDB := stdlib.OpenDBFromPool(connPool)
	db := sqlx.NewDb(nativeDB, "pgx")

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	m, err := migrate.NewWithDatabaseInstance("embed://", config.Database.Name, driver)
	if err != nil {
		log.Fatalln(err)
	}

	if err := m.Up(); errors.Is(err, migrate.ErrNoChange) {
		log.Println(err)
	} else if err != nil {
		log.Fatalln(err)
	}

	return db, nil
}
