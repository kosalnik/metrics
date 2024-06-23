package backup

import (
	"context"
	"encoding/json"
	"os"

	"github.com/kosalnik/metrics/internal/infra/storage"
)

// Dump - Тип выполняет сохранение содержимого хранилища на диск.
type Dump struct {
	storage storage.Dumper
	path    string
}

// NewDump возвращает тип Dump.
// На вход ожидаются Storage и абсолютный путь до файла, в который нужно сохранять бекап
func NewDump(storage storage.Dumper, path string) *Dump {
	return &Dump{
		storage: storage,
		path:    path,
	}
}

// Store - вызовом этого метода данные из Storage сохраняются на диск
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
