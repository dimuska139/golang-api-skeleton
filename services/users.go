package services

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	apiErrors "github.com/dimuska139/golang-api-skeleton/api_errors"
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

func (a *UsersService) List() ([]dto.UserDTO, error) {
	// Using raw sql
	rows, err := a.Database.Queryx("SELECT * FROM users ORDER BY id")

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

func (a *UsersService) GetByEmail(email string) (*dto.UserDTO, error) {
	userDTO := dto.UserDTO{}
	q, _, _ := sq.Select("id", "name", "email", "hash").
		From("users").
		Where("email = ?").
		Limit(1).
		ToSql()

	if err := a.Database.Get(&userDTO, q, email); err != nil {
		return nil, &apiErrors.NotFoundError{S: "User not found"}
	}

	return &userDTO, nil
}

func (a *UsersService) GetByToken(refreshToken string) (*dto.UserDTO, error) {
	userDTO := dto.UserDTO{}
	q, _, _ := sq.Select("users.id", "users.name", "users.email", "users.hash").
		From("users").
		Join("jwt_refresh tokens ON users.id=tokens.user_id").
		Where("tokens.refresh_token = ?").
		Limit(1).
		ToSql()


	if err := a.Database.Get(&userDTO, q, refreshToken); err != nil {
		return nil, &apiErrors.NotFoundError{S: "User not found"}
	}

	return &userDTO, nil
}

func (a *UsersService) Create(email string, name string, hash string) (*dto.UserDTO, error) {
	n := sql.NullString{}
	if name != "" {
		n = sql.NullString{
			String: name,
			Valid:  true,
		}
	}

	// Using Squirrel (https://github.com/Masterminds/squirrel)
	query := sq.Insert("users").
		Columns("email", "name", "hash").
		Values(email, n, hash).
		Suffix("RETURNING \"id\"").
		RunWith(a.Database.DB).
		PlaceholderFormat(sq.Dollar)

	var resName *string = nil
	if len(name) != 0 {
		resName = &name
	}

	res := dto.UserDTO{
		Email: email,
		Name:  resName,
	}
	if err := query.QueryRow().Scan(&res.ID); err != nil {
		return nil, errors.Wrap(err, "Insert user to database")
	}
	return &res, nil
}
