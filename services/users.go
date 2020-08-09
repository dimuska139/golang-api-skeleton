package services

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/dimuska139/golang-api-skeleton/dto"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type UsersService struct {
	Database *sqlx.DB
}

func NewUsersService(database *sqlx.DB) *UsersService {
	return &UsersService{Database: database}
}

func (a *UsersService) CountTotal() (int, error) {
	row := a.Database.QueryRow("SELECT COUNT(*) FROM users")
	var res int
	if err := row.Scan(&res); err != nil {
		return 0, errors.Wrap(err, "Query error")
	}
	return res, nil
}

func (a *UsersService) Create(email string, name string) (*dto.UserDTO, error) {
	// Using Squirrel (https://github.com/Masterminds/squirrel)
	query := sq.Insert("users").
		Columns("email", "name").
		Values(email, name).
		Suffix("RETURNING \"id\"").
		RunWith(a.Database.DB).
		PlaceholderFormat(sq.Dollar)

	res := dto.UserDTO{
		Email: email,
		Name:  name,
	}

	if err := query.QueryRow().Scan(&res.ID); err != nil {
		return nil, errors.Wrap(err, "Insert user to database")
	}
	return &res, nil
}

func (a *UsersService) List() ([]dto.UserDTO, error) {
	// Using raw sql
	rows, err := a.Database.Queryx("SELECT * FROM users")

	if err != nil {
		return nil, errors.Wrap(err, "Getting users from database")
	}

	var items = make([]dto.UserDTO, 0)
	for rows.Next() {
		var item dto.UserDTO
		rows.StructScan(&item)
		items = append(items, item)
	}

	return items, nil
}
