package token

import (
	"context"
	"errors"
	"github.com/JAbduvohidov/jwt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Service struct {
	secret jwt.Secret
	pool   *pgxpool.Pool
}

func NewService(secret jwt.Secret, pool *pgxpool.Pool) *Service {
	return &Service{secret: secret, pool: pool}
}

var (
	ROLE_USER  = "USER"
	ROLE_ADMIN = "ADMIN"
)

type Payload struct {
	Id    int64  `json:"id"`
	Login string `json:"login"`
	Role  string `json:"role"`
	Exp   int64  `json:"exp"`
}

type RequestDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type ResponseDTO struct {
	Token string `json:"token"`
}

var ErrInvalidLogin = errors.New("invalid password")
var ErrInvalidPassword = errors.New("invalid password")

func (s *Service) Generate(ctx context.Context, request *RequestDTO) (response ResponseDTO, err error) {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return response, err
	}

	row := conn.QueryRow(ctx, getUserByLoginDML, request.Login)

	var userPayload Payload
	var hPassword string
	err = row.Scan(&userPayload.Id, &userPayload.Login, &hPassword, &userPayload.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = ErrInvalidLogin
		}
		return
	}

	userPayload.Exp = time.Now().Add(time.Hour * 10).Unix()

	err = bcrypt.CompareHashAndPassword([]byte(hPassword), []byte(request.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		err = ErrInvalidPassword
		return
	}

	response.Token, err = jwt.Encode(userPayload, s.secret)
	return
}
