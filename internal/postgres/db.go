package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

func NewConn(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var ok bool
	row := db.QueryRowContext(ctx, "SELECT true AS ok")
	if err := row.Scan(&ok); err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("DB Connection timeout exceed")
	}
	return db, nil
}
