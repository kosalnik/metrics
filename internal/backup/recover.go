package backup

import (
	"context"
	"encoding/json"
	"os"

	"github.com/kosalnik/metrics/internal/log"
	"github.com/kosalnik/metrics/internal/storage"
)

type Recover struct {
	recover storage.BatchInserter
	path    string
}

// NewRecover - создаст инстанс типа Recover, которую можно использовать для восстановления Storage из бекапа.
// На вход подаётся объект реализующий интерфейс storage.BatchInserter и путь к файлу из которого нужно восстанавливать.
func NewRecover(storage storage.BatchInserter, path string) *Recover {
	return &Recover{
		recover: storage,
		path:    path,
	}
}

// Recover - Восстановить Storage из бекапа.
func (m *Recover) Recover(ctx context.Context) error {
	if m == nil || m.path == "" {
		log.Info().Msg("Recover skipped. No Path or Disabled")

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
	if err := m.recover.UpsertAll(ctx, b.Data); err != nil {
		return err
	}

	return nil
}
