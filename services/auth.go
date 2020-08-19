package services

import (
	"crypto/rand"
	"encoding/hex"
	sq "github.com/Masterminds/squirrel"
	"github.com/dgrijalva/jwt-go"
	apiErrors "github.com/dimuska139/golang-api-skeleton/api_errors"
	"github.com/dimuska139/golang-api-skeleton/config"
	"github.com/dimuska139/golang-api-skeleton/dto"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService struct {
	Database     *sqlx.DB
	Config       *config.Config
	UsersService *UsersService
}

func NewAuthService(database *sqlx.DB, config *config.Config, usersService *UsersService) *AuthService {
	return &AuthService{Database: database, Config: config, UsersService: usersService}
}

func (a *AuthService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (a *AuthService) saveRefreshToken(userId int, jwtRefresh string, expiresAt time.Time) error {
	query := sq.Insert("jwt_refresh").
		Columns("user_id", "refresh_token", "expires_at").
		Values(userId, jwtRefresh, expiresAt).
		RunWith(a.Database.DB).
		PlaceholderFormat(sq.Dollar)

	_, err := query.Query()

	if err != nil {
		return errors.Wrap(err, "Insert refresh token to database")
	}

	return nil
}

func (a *AuthService) GenerateJwtPair(userDTO dto.UserDTO) (*dto.JwtPairDTO, error) {
	dt := time.Now()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &dto.ClaimsDTO{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: dt.Add(time.Duration(a.Config.Jwt.AccessTokenLifetime) * time.Second).Unix(),
			IssuedAt:  dt.Unix(),
		},
		User: userDTO,
	})

	access, err := accessToken.SignedString([]byte(a.Config.Jwt.Secret))
	if err != nil {
		return nil, errors.Wrap(err, "Access token signing")
	}

	refreshTokenB := make([]byte, 32)
	if _, err := rand.Read(refreshTokenB); err != nil {
		return nil, errors.Wrap(err, "Generate refresh token")
	}

	refresh := hex.EncodeToString(refreshTokenB)

	if err := a.saveRefreshToken(userDTO.ID, refresh, dt.Add(time.Duration(a.Config.Jwt.RefreshTokenLifetime)*time.Second)); err != nil {
		return nil, errors.Wrap(err, "Save refresh")
	}

	return &dto.JwtPairDTO{
		Access:  access,
		Refresh: refresh,
	}, nil
}

func (a *AuthService) InvalidateOldToken(refreshToken string) error {
	// Appends 30 seconds handle multiple async requests from one client app. Without it the first request generates new
	// refresh token and invalidates previous. So another requests fails because their token expired.
	_, err := a.Database.Exec(`UPDATE jwt_refresh SET expires_at=now() + INTERVAL '30 seconds'
		WHERE refresh_token=$1`, refreshToken)
	return err
}

func (a *AuthService) Login(email string, password string) (*dto.UserDTO, error) {
	userDTO, err := a.UsersService.GetByEmail(email)
	if err != nil {
		return nil, &apiErrors.NotFoundError{S: "User not found"}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userDTO.Hash), []byte(password)); err != nil {
		return nil, &apiErrors.NotFoundError{S: "Invalid password"}
	}

	return userDTO, nil
}

func (a *AuthService) Registration(email string, name string, password string) (*dto.UserDTO, error) {
	hash, err := a.hashPassword(password)
	if err != nil {
		return nil, errors.Wrap(err, "Hash password")
	}

	userDTO, err := a.UsersService.Create(email, name, hash)

	if err != nil {
		return nil, errors.Wrap(err, "Create new user")
	}

	return userDTO, nil
}
