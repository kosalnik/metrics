package backup

import (
	"context"
	"encoding/json"
	"os"

	"github.com/kosalnik/metrics/internal/infra/logger"
	"github.com/kosalnik/metrics/internal/infra/storage"
)

type Recoverer interface {
	Recover(ctx context.Context) error
}

type Recover struct {
	storage storage.Storage
	path    string
}

func NewRecover(storage storage.Storage, path string) *Recover {
	return &Recover{
		storage: storage,
		path:    path,
	}
}

func (m *Recover) Recover(ctx context.Context) error {
	if m == nil || m.path == "" {
		logger.Logger.Info("Recover skipped. No Path or Disabled")

		return nil
	}
	d, err := os.ReadFile(m.path)
	if err != nil {
		return err
	}
	var b Backup
	if err := json.Unmarshal(d, &b); err != nil {
		return err
	}
	if len(b.Data) == 0 {
		return nil
	}
	if err := m.storage.UpsertAll(ctx, b.Data); err != nil {
		return err
	}

	return nil
}
