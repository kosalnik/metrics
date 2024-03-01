package entity

import "testing"

func TestCounterValue_Inc(t *testing.T) {
	tests := []struct {
		name  string
		init  int64
		value int64
		want  int64
	}{
		{name: "add", init: 1, value: 2, want: 3},
		{name: "del", init: 1, value: -2, want: -1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := CounterValue{Name: "a", Value: test.init}
			m.Inc(test.value)
			if m.Value != test.want {
				t.Errorf("expected %v, got %v", test.want, m.Value)
			}
		})
	}
}
