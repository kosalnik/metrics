package storage

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestMemStorage_GetCounter(t *testing.T) {
	type storageState struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	tests := []struct {
		name         string
		storageState storageState
		metricName   string
		want         int64
	}{
		{
			name: "empty storage",
			storageState: storageState{
				gauge:   map[string]float64{},
				counter: map[string]int64{},
			},
			metricName: "testCounter",
			want:       0,
		},
		{
			name: "no gauge, counter exists",
			storageState: storageState{
				gauge: map[string]float64{},
				counter: map[string]int64{
					"testCounter": 3,
				},
			},
			metricName: "testCounter",
			want:       3,
		},
		{
			name: "gauge exists, no counter",
			storageState: storageState{
				gauge: map[string]float64{
					"testCounter": 3,
				},
				counter: map[string]int64{},
			},
			metricName: "testCounter",
			want:       0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				gauge:   tt.storageState.gauge,
				counter: tt.storageState.counter,
			}
			if got := m.GetCounter(tt.metricName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
	type storageState struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	tests := []struct {
		name         string
		storageState storageState
		metricName   string
		want         float64
	}{
		{
			name: "empty storage",
			storageState: storageState{
				gauge:   map[string]float64{},
				counter: map[string]int64{},
			},
			metricName: "testCounter",
			want:       0,
		},
		{
			name: "no gauge, counter exists",
			storageState: storageState{
				gauge: map[string]float64{},
				counter: map[string]int64{
					"testCounter": 3,
				},
			},
			metricName: "testCounter",
			want:       0,
		},
		{
			name: "gauge exists, no counter",
			storageState: storageState{
				gauge: map[string]float64{
					"testCounter": 3,
				},
				counter: map[string]int64{},
			},
			metricName: "testCounter",
			want:       3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				gauge:   tt.storageState.gauge,
				counter: tt.storageState.counter,
			}
			if got := m.GetGauge(tt.metricName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGauge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_IncCounter(t *testing.T) {
	type storageState struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	type metric struct {
		name  string
		value int64
	}
	tests := []struct {
		name         string
		storageState storageState
		metric       metric
		want         int64
	}{
		{
			name: "empty storage",
			storageState: storageState{
				gauge:   map[string]float64{},
				counter: map[string]int64{},
			},
			metric: metric{name: "test", value: 3},
			want:   3,
		},
		{
			name: "no gauge, counter exists",
			storageState: storageState{
				gauge: map[string]float64{},
				counter: map[string]int64{
					"test": 2,
				},
			},
			metric: metric{name: "test", value: 3},
			want:   5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				gauge:   tt.storageState.gauge,
				counter: tt.storageState.counter,
			}
			m.IncCounter(tt.metric.name, tt.metric.value)
			actual := m.GetCounter(tt.metric.name)
			assert.Equal(t, tt.want, actual)
		})
	}
}

func TestMemStorage_SetGauge(t *testing.T) {
	m := NewStorage()
	assert.NotNil(t, m)
	m.SetGauge("test", 1)
	v := m.GetGauge("test")
	assert.NotNil(t, v)
	assert.Equal(t, 1.0, v)
}

func TestNewStorage(t *testing.T) {
	m := NewStorage()
	assert.NotNil(t, m)
}
