package models

import "testing"

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
