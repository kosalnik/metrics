package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kosalnik/metrics/internal/config"
	"github.com/kosalnik/metrics/internal/models"
	"github.com/sirupsen/logrus"
)

var schemaGaugeSQL = `CREATE TABLE IF NOT EXISTS gauge(
    	id VARCHAR(200) PRIMARY KEY,
    	value double precision not null
    )`

var schemaCounterSQL = `CREATE TABLE IF NOT EXISTS counter(
    	id VARCHAR(200) PRIMARY KEY,
    	value double precision not null
    )`

type DBStorage struct {
	db *sql.DB
}

func NewDB(ctx context.Context, cfg config.DB) (*DBStorage, error) {
	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		return nil, err
	}

	if _, err := db.ExecContext(ctx, schemaCounterSQL); err != nil {
		return nil, err
	}

	if _, err := db.ExecContext(ctx, schemaGaugeSQL); err != nil {
		return nil, err
	}

	return &DBStorage{db}, nil
}

func (d DBStorage) GetGauge(ctx context.Context, name string) (float64, bool, error) {
	r := d.db.QueryRowContext(ctx, "SELECT value FROM gauge WHERE id = $1", name)
	var v float64
	if err := r.Scan(&v); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, false, nil
		}
		logrus.WithError(err).Error("fail db request. get gauge")
		return 0, false, err
	}
	return v, true, nil
}

func (d DBStorage) SetGauge(ctx context.Context, name string, value float64) (float64, error) {
	s := "INSERT INTO gauge (id, value) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET value = $2"
	_, err := d.db.ExecContext(ctx, s, name, value)
	if err != nil {
		logrus.WithError(err).Error("fail db request. set gauge")
		return 0, err
	}
	v, _, err := d.GetGauge(ctx, name)
	return v, err
}

func (d DBStorage) GetCounter(ctx context.Context, name string) (int64, bool, error) {
	r := d.db.QueryRowContext(ctx, "SELECT value FROM counter WHERE id = $1", name)
	var v int64
	if err := r.Scan(&v); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, false, nil
		}
		logrus.WithError(err).Error("fail db request. get counter")
		return 0, false, err
	}
	return v, true, nil
}

func (d DBStorage) IncCounter(ctx context.Context, name string, value int64) (int64, error) {
	s := "INSERT INTO counter (id, value) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET value = counter.value + $2"
	_, err := d.db.ExecContext(ctx, s, name, value)
	if err != nil {
		logrus.WithError(err).Error("fail db request. inc counter")
		return 0, err
	}
	v, _, err := d.GetCounter(ctx, name)
	return v, err
}

func (d DBStorage) GetAll(ctx context.Context) ([]models.Metrics, error) {
	g, err := d.getAllGauge(ctx)
	if err != nil {
		return nil, err
	}
	c, err := d.getAllCounter(ctx)
	if err != nil {
		return nil, err
	}

	return append(g, c...), nil
}

func (d DBStorage) getAllGauge(ctx context.Context) ([]models.Metrics, error) {
	var res []models.Metrics
	rows, err := d.db.QueryContext(ctx, "SELECT id, value FROM gauge ORDER BY id")
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		m := models.Metrics{MType: models.MGauge}
		if err := rows.Scan(&m.ID, &m.Value); err != nil {
			return nil, err
		}
		res = append(res, m)
	}

	return res, nil
}

func (d DBStorage) getAllCounter(ctx context.Context) ([]models.Metrics, error) {
	var res []models.Metrics
	rows, err := d.db.QueryContext(ctx, "SELECT id, value FROM counter ORDER BY id")
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		m := models.Metrics{MType: models.MCounter}
		if err := rows.Scan(&m.ID, &m.Delta); err != nil {
			return nil, err
		}
		res = append(res, m)
	}

	return res, nil
}

func (d DBStorage) Close() error {
	return d.db.Close()
}

func (d DBStorage) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

func (d DBStorage) Store(ctx context.Context, path string) error {
	return nil
}

func (d DBStorage) Recover(ctx context.Context, path string) error {
	return nil
}
