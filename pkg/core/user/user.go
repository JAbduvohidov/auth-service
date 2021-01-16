package user

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

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
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
	Avatar   string `json:"avatar"`
}

func (s *Service) GetProfile(ctx context.Context) (response RequestDTO, err error) {
	err = s.pool.QueryRow(ctx, getUserDML, ctx.Value("id")).Scan(&response.Id, &response.Name, &response.Surname, &response.Login, &response.Avatar)
	if err != nil {
		return RequestDTO{}, fmt.Errorf("unable to get user info: %w", err)
	}
	return
}

func (s *Service) EditProfile(ctx context.Context, user *RequestDTO) (err error) {
	_, err = s.pool.Exec(ctx, updateUserProfileDML, user.Name, user.Surname, ctx.Value("id"))
	if err != nil {
		return fmt.Errorf("unable to update user info: %w", err)
	}

	if user.Password == "" {
		return
	}

	hPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	_, err = s.pool.Exec(ctx, updateUserPasswordDML, hPassword, ctx.Value("id"))
	if err != nil {
		return fmt.Errorf("unable to update user password: %w", err)
	}
	return
}

func (s *Service) AddUser(ctx context.Context, request RequestDTO) (err error) {
	hPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	_, err = s.pool.Exec(ctx, addUserDML, request.Name, request.Surname, request.Login, hPassword, request.Avatar)
	if err != nil {
		return fmt.Errorf("unable to add user: %w", err)
	}
	return nil
}

func (s *Service) GetUsers(ctx context.Context) (users []*RequestDTO, err error) {
	rows, err := s.pool.Query(ctx, getUsersDML)
	if err != nil {
		return users, fmt.Errorf("unable to get users: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		user := RequestDTO{}
		err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.Surname,
			&user.Login,
			&user.Role,
			&user.Avatar,
		)
		if err != nil {
			return users, fmt.Errorf("unable to scan user: %w", err)
		}
		users = append(users, &user)
	}
	return
}

func (s *Service) DeleteUser(ctx context.Context, id int) (err error) {
	_, err = s.pool.Exec(ctx, deleteUserDML, id)
	if err != nil {
		return fmt.Errorf("unable to delete user: %w", err)
	}
	return
}

func (s *Service) UpgradeUser(ctx context.Context, role string, id int) (err error) {
	_, err = s.pool.Exec(ctx, upgradeUserDML, role, id)
	if err != nil {
		return fmt.Errorf("unable to upgrade user: %w", err)
	}
	return
}
