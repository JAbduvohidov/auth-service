package services

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

func InitDB(dsn string) (err error) {
	conn, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("unable to connect to pool with dsn: %s, %w", dsn, err)
	}
	_, err = conn.Query(context.Background(), usersDDL)
	if err != nil {
		return fmt.Errorf("unable to create table: %w", err)
	}

	_, err = conn.Query(context.Background(), moderatorDML)
	if err != nil {
		return fmt.Errorf("unable to add moderator: %w", err)
	}
	return nil
}
