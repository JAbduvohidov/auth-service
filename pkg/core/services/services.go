package services

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func InitDB(dsn string, mPass string) (err error) {
	conn, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("unable to connect to pool with dsn: %s, %w", dsn, err)
	}
	_, err = conn.Query(context.Background(), usersDDL)
	if err != nil {
		return fmt.Errorf("unable to create table: %w", err)
	}

	time.Sleep(time.Second)

	hPassword, err := bcrypt.GenerateFromPassword([]byte(mPass), bcrypt.DefaultCost)
	_, err = conn.Query(context.Background(), moderatorDML, hPassword)
	if err != nil {
		return fmt.Errorf("unable to add moderator: %w", err)
	}
	return nil
}
