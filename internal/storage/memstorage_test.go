package storage

import (
	"github.com/kosalnik/metrics/internal/entity"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestMemStorage_GetCounter(t *testing.T) {
	type storageState struct {
		gauge   map[string]*entity.GaugeValue
		counter map[string]*entity.CounterValue
	}
	tests := []struct {
		name         string
		storageState storageState
		metricName   string
		want         *entity.CounterValue
	}{
		{
			name: "empty storage",
			storageState: storageState{
				gauge:   map[string]*entity.GaugeValue{},
				counter: map[string]*entity.CounterValue{},
			},
			metricName: "testCounter",
			want:       nil,
		},
		{
			name: "no gauge, counter exists",
			storageState: storageState{
				gauge: map[string]*entity.GaugeValue{},
				counter: map[string]*entity.CounterValue{
					"testCounter": {Name: "testCounter", Value: 3},
				},
			},
			metricName: "testCounter",
			want:       &entity.CounterValue{Name: "testCounter", Value: 3},
		},
		{
			name: "gauge exists, no counter",
			storageState: storageState{
				gauge: map[string]*entity.GaugeValue{
					"testCounter": {Name: "testCounter", Value: 3},
				},
				counter: map[string]*entity.CounterValue{},
			},
			metricName: "testCounter",
			want:       nil,
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
		gauge   map[string]*entity.GaugeValue
		counter map[string]*entity.CounterValue
	}
	tests := []struct {
		name         string
		storageState storageState
		metricName   string
		want         *entity.GaugeValue
	}{
		{
			name: "empty storage",
			storageState: storageState{
				gauge:   map[string]*entity.GaugeValue{},
				counter: map[string]*entity.CounterValue{},
			},
			metricName: "testCounter",
			want:       nil,
		},
		{
			name: "no gauge, counter exists",
			storageState: storageState{
				gauge: map[string]*entity.GaugeValue{},
				counter: map[string]*entity.CounterValue{
					"testCounter": {Name: "testCounter", Value: 3},
				},
			},
			metricName: "testCounter",
			want:       nil,
		},
		{
			name: "gauge exists, no counter",
			storageState: storageState{
				gauge: map[string]*entity.GaugeValue{
					"testCounter": {Name: "testCounter", Value: 3},
				},
				counter: map[string]*entity.CounterValue{},
			},
			metricName: "testCounter",
			want:       &entity.GaugeValue{Name: "testCounter", Value: 3},
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
		gauge   map[string]*entity.GaugeValue
		counter map[string]*entity.CounterValue
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
				gauge:   map[string]*entity.GaugeValue{},
				counter: map[string]*entity.CounterValue{},
			},
			metric: metric{name: "test", value: 3},
			want:   3,
		},
		{
			name: "no gauge, counter exists",
			storageState: storageState{
				gauge: map[string]*entity.GaugeValue{},
				counter: map[string]*entity.CounterValue{
					"test": {Name: "testCounter", Value: 2},
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
			assert.NotNil(t, actual)
			assert.Equal(t, tt.want, actual.Value)
		})
	}
}

func TestMemStorage_SetGauge(t *testing.T) {
	m := NewStorage()
	assert.NotNil(t, m)
	m.SetGauge("test", 1)
	v := m.GetGauge("test")
	assert.NotNil(t, v)
	assert.Equal(t, "test", v.Name)
	assert.Equal(t, 1.0, v.Value)
}

func TestNewStorage(t *testing.T) {
	m := NewStorage()
	assert.NotNil(t, m)
}
