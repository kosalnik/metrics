package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetrics_String(t *testing.T) {
	type fields struct {
		MType MType
		ID    string
		Delta int64
		Value float64
	}
	tests := map[string]struct {
		want   string
		fields fields
	}{
		"counter": {
			fields: fields{ID: "zxc", MType: "counter", Delta: 10},
			want:   "zxc = 10",
		},
		"gauge ceil": {
			fields: fields{ID: "zxc", MType: "gauge", Value: 10},
			want:   "zxc = 10",
		},
		"gauge float": {
			fields: fields{ID: "zxc", MType: "gauge", Value: 3.1415},
			want:   "zxc = 3.1415",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			m := &Metrics{
				ID:    tt.fields.ID,
				MType: tt.fields.MType,
				Delta: tt.fields.Delta,
				Value: tt.fields.Value,
			}
			if got := m.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetrics_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		obj  *Metrics
		want string
	}{
		{name: "Gauge", obj: &Metrics{ID: "aa", MType: MGauge, Value: 2.14}, want: `{"id":"aa","type":"gauge","value":2.14}`},
		{name: "Counter", obj: &Metrics{ID: "bb", MType: MCounter, Delta: 3}, want: `{"id":"bb","type":"counter","delta":3}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.obj.MarshalJSON()
			require.NoError(t, err)
			assert.JSONEq(t, tt.want, string(got))
		})
	}
}
