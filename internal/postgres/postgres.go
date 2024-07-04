// Package postgres содержит реализацию репозиториев к БД postgres.
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/kosalnik/metrics/internal/logger"
	"github.com/kosalnik/metrics/internal/models"
)

var schemaGaugeSQL = `CREATE TABLE IF NOT EXISTS gauge(
    	id VARCHAR(200) PRIMARY KEY,
    	value double precision not null
    )`

var schemaCounterSQL = `CREATE TABLE IF NOT EXISTS counter(
    	id VARCHAR(200) PRIMARY KEY,
    	value bigint not null
    )`

type DBStorage struct {
	updatedAt time.Time
	db        *sql.DB
	mu        sync.Mutex
}

func NewDBStorage(db *sql.DB) (*DBStorage, error) {
	return &DBStorage{mu: sync.Mutex{}, updatedAt: time.Now(), db: db}, nil
}

func (d *DBStorage) InitTables(ctx context.Context) error {
	if _, err := d.db.ExecContext(ctx, schemaCounterSQL); err != nil {
		return err
	}
	if _, err := d.db.ExecContext(ctx, schemaGaugeSQL); err != nil {
		return err
	}
	return nil
}

func (d *DBStorage) GetGauge(ctx context.Context, name string) (*models.Metrics, error) {
	r := d.db.QueryRowContext(ctx, "SELECT value FROM gauge WHERE id = $1", name)
	var v float64
	if err := r.Scan(&v); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		logger.Logger.WithError(err).Error("fail db request. get gauge")
		return nil, err
	}
	return &models.Metrics{ID: name, MType: models.MGauge, Value: v}, nil
}

func (d *DBStorage) SetGauge(ctx context.Context, name string, value float64) (*models.Metrics, error) {
	s := "INSERT INTO gauge (id, value) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET value = $2"
	_, err := d.db.ExecContext(ctx, s, name, value)
	if err != nil {
		logger.Logger.WithError(err).Error("fail db request. set gauge")

		return nil, err
	}
	d.setUpdatedAt()
	v, err := d.GetGauge(ctx, name)

	return v, err
}

func (d *DBStorage) GetCounter(ctx context.Context, name string) (*models.Metrics, error) {
	r := d.db.QueryRowContext(ctx, "SELECT value FROM counter WHERE id = $1", name)
	var v int64
	if err := r.Scan(&v); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		logger.Logger.WithError(err).Error("fail db request. get counter")

		return nil, err
	}
	return &models.Metrics{ID: name, MType: models.MCounter, Delta: v}, nil
}

func (d *DBStorage) IncCounter(ctx context.Context, name string, value int64) (*models.Metrics, error) {
	s := "INSERT INTO counter (id, value) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET value = counter.value + $2"
	_, err := d.db.ExecContext(ctx, s, name, value)
	if err != nil {
		logger.Logger.WithError(err).Error("fail db request. inc counter")
		return nil, err
	}
	d.setUpdatedAt()
	v, err := d.GetCounter(ctx, name)

	return v, err
}

func (d *DBStorage) inTransaction(ctx context.Context, fn func(tr *sql.Tx) error) error {
	tr, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	err = fn(tr)
	if err != nil {
		err := tr.Rollback()
		if err != nil {
			return fmt.Errorf("fail rollback transaction: %w", err)
		}
		return err
	}
	return nil
}

func (d *DBStorage) UpsertAll(ctx context.Context, list []models.Metrics) (err error) {
	return d.inTransaction(ctx, func(tr *sql.Tx) error {
		incCounterSt, err := tr.PrepareContext(
			ctx,
			"INSERT INTO counter (id, value) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET value = counter.value + $2",
		)
		if err != nil {
			return fmt.Errorf("fail upsert: %w", err)
		}
		defer incCounterSt.Close()
		setGaugeSt, err := tr.PrepareContext(
			ctx,
			"INSERT INTO gauge (id, value) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET value = $2",
		)
		if err != nil {
			return fmt.Errorf("fail upsert: %w", err)
		}
		defer setGaugeSt.Close()
		logger.Logger.WithField("list", list).Info("upsertAll")
		for _, v := range list {
			switch v.MType {
			case models.MGauge:
				if _, err := setGaugeSt.ExecContext(ctx, v.ID, v.Value); err != nil {
					return err
				}
				continue
			case models.MCounter:
				if _, err := incCounterSt.ExecContext(ctx, v.ID, v.Delta); err != nil {
					return err
				}
			}
		}
		if err := tr.Commit(); err != nil {
			return fmt.Errorf("fail commit transaction: %w", err)
		}
		d.setUpdatedAt()
		return nil
	})
}

func (d *DBStorage) GetAll(ctx context.Context) ([]models.Metrics, error) {
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

func (d *DBStorage) getAllGauge(ctx context.Context) ([]models.Metrics, error) {
	var res []models.Metrics
	rows, err := d.db.QueryContext(ctx, "SELECT id, value FROM gauge ORDER BY id")
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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (d *DBStorage) getAllCounter(ctx context.Context) ([]models.Metrics, error) {
	var res []models.Metrics
	rows, err := d.db.QueryContext(ctx, "SELECT id, value FROM counter ORDER BY id")
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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (d *DBStorage) Close() error {
	return d.db.Close()
}

func (d *DBStorage) UpdatedAt() time.Time {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.updatedAt
}

func (d *DBStorage) setUpdatedAt() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.updatedAt = time.Now()
}

func (d *DBStorage) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}
