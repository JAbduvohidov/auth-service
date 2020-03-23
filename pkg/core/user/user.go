package user

import (
	"auth-service/pkg/core/token"
	"context"
	"errors"
	"fmt"
	"github.com/JAbduvohidov/mux/middleware/jwt"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var ErrStructType = errors.New("incorrect struct type")

type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

type ResponseDTO struct {
	Id    int64  `json:"id"`
	Login string `json:"login"`
	Role  string `json:"role"`
}

type RequestDTO struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
}

func (s *Service) Profile(ctx context.Context) (response ResponseDTO, err error) {
	auth, ok := jwt.FromContext(ctx).(*token.Payload)
	if !ok {
		return ResponseDTO{}, ErrStructType
	}

	return ResponseDTO{
		Id:    auth.Id,
		Login: auth.Login,
		Role:  auth.Role,
	}, nil
}

func (s *Service) AddUser(ctx context.Context, request RequestDTO) (err error) {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("unable to acquire ctx: %w", err)
	}

	hPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	_, err = conn.Exec(ctx, addUserDML, request.Name, request.Surname, request.Login, hPassword, request.Avatar)
	if err != nil {
		return fmt.Errorf("unable to add user: %w", err)
	}
	return nil
}
