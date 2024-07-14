package backup

import (
	"context"
	"time"

	"github.com/kosalnik/metrics/internal/log"
	"github.com/kosalnik/metrics/internal/models"
	"github.com/kosalnik/metrics/internal/storage"
)

type Dumper interface {
	Store(ctx context.Context) error
}

type Recoverer interface {
	Recover(ctx context.Context) error
}

type Storage interface {
	storage.Dumper
	storage.BatchInserter
	storage.UpdateAwarer
}

type Backup struct {
	Data []models.Metrics
}

// BackupManager - структура управляющая бекапом и восстановлением из бекапа.
type BackupManager struct {
	dump           Dumper
	recover        Recoverer
	storage        Storage
	lastBackup     time.Time
	backupInterval time.Duration
}

// NewBackupManager - создаёт BackupManager с конфигурацией config.Backup.
// На вход принимает хранилище с интерфейсом Storage и настройки бекапа.
func NewBackupManager(s Storage, cfg Config) (*BackupManager, error) {
	if cfg.FileStoragePath == "" {
		return nil, nil
	}
	var d *Dump
	var r *Recover
	if cfg.StoreInterval > 0 {
		d = NewDump(s, cfg.FileStoragePath)
	}
	if cfg.Restore {
		r = NewRecover(s, cfg.FileStoragePath)
	}
	if d == nil && r == nil {
		return nil, nil
	}
	return &BackupManager{
		dump:           d,
		recover:        r,
		storage:        s,
		backupInterval: time.Duration(cfg.StoreInterval) * time.Second,
		lastBackup:     time.Now(),
	}, nil
}

func (m *BackupManager) BackupLoop(ctx context.Context) {
	if m == nil || m.dump == nil || m.backupInterval == 0 {
		log.Info().Msg("schedule backup skipped")
		return
	}

	tick := time.NewTicker(m.backupInterval)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("backup loop: context done")
			return
		case <-tick.C:
			if m.lastBackup.Equal(m.storage.UpdatedAt()) {
				log.Debug().Msg("backup loop: no changes, skip backup")
				continue
			}
			log.Info().Msg("backup loop: store")
			if err := m.dump.Store(ctx); err != nil {
				log.Error().Err(err).Msg("Fail backup")
			}
		}
	}
}

// Recover - восстановить данные из бекапа.
func (m *BackupManager) Recover(ctx context.Context) error {
	if m == nil || m.recover == nil {
		log.Info().Msg("recover skipped")

		return nil
	}

	log.Info().Msg("recover start")
	return m.recover.Recover(ctx)
}
