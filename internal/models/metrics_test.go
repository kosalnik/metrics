package models

import "testing"

func TestMetrics_String(t *testing.T) {
	type fields struct {
		ID    string
		MType MType
		Delta *int64
		Value *float64
	}
	var tenInt int64 = 10
	var tenFloat float64 = 10
	var pi float64 = 3.1415

	tests := map[string]struct {
		fields fields
		want   string
	}{
		"counter": {
			fields: fields{ID: "zxc", MType: "counter", Delta: &tenInt},
			want:   "zxc = 10",
		},
		"gauge ceil": {
			fields: fields{ID: "zxc", MType: "gauge", Value: &tenFloat},
			want:   "zxc = 10",
		},
		"gauge float": {
			fields: fields{ID: "zxc", MType: "gauge", Value: &pi},
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
