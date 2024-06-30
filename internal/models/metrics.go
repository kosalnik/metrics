// Package models содержит определения моделей, используемых приложением.
package models

import (
	"encoding/json"
	"fmt"
)

type Metrics struct {
	ID    string  `json:"id"`              // имя метрики
	MType MType   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *Metrics) MarshalJSON() ([]byte, error) {
	switch m.MType {
	case MGauge:
		return json.Marshal(
			struct {
				ID    string  `json:"id"`
				MType MType   `json:"type"`
				Value float64 `json:"value"`
			}{m.ID, m.MType, m.Value},
		)
	case MCounter:
		return json.Marshal(
			struct {
				ID    string `json:"id"`
				MType MType  `json:"type"`
				Delta int64  `json:"delta"`
			}{m.ID, m.MType, m.Delta},
		)
	}

	return json.Marshal(m)
}

type MType string

const (
	MGauge   MType = "gauge"
	MCounter MType = "counter"
)

func (m *Metrics) String() string {
	if m.MType == MCounter {
		return fmt.Sprintf("%s = %d", m.ID, m.Delta)
	}

	return fmt.Sprintf("%s = %g", m.ID, m.Value)
}
