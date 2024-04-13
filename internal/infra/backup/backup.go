package backup

import (
	"context"
	"encoding/json"
	"os"

	"github.com/kosalnik/metrics/internal/infra/storage"
)

type Dumper interface {
	Store(ctx context.Context) error
}

type Dump struct {
	storage storage.Storage
	path    string
}

func NewDump(storage storage.Storage, path string) *Dump {
	return &Dump{
		storage: storage,
		path:    path,
	}
}

func (m *Dump) Store(ctx context.Context) error {
	if m.path == "" {
		return nil
	}
	f, err := os.CreateTemp(os.TempDir(), "backup")
	if err != nil {
		return err
	}
	savePath := f.Name()
	b, err := m.storage.GetAll(ctx)
	if err != nil {
		return err
	}
	d, err := json.Marshal(Backup{Data: b})
	if err != nil {
		return err
	}
	if _, err := f.Write(d); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if err := os.Rename(savePath, m.path); err != nil {
		return err
	}
	return nil
}
