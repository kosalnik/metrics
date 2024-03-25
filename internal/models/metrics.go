package models

import "fmt"

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType MType    `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type MType string

const (
	MGauge   MType = "gauge"
	MCounter MType = "counter"
)

func (m *Metrics) String() string {
	if m.MType == MCounter {
		return fmt.Sprintf("%s = %d", m.ID, *m.Delta)
	}

	return fmt.Sprintf("%s = %g", m.ID, *m.Value)
}
