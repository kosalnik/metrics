package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/kosalnik/metrics/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		wantOk       bool
	}{
		{
			name: "empty storage",
			storageState: storageState{
				gauge:   map[string]float64{},
				counter: map[string]int64{},
			},
			metricName: "testCounter",
			want:       0,
			wantOk:     false,
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
			wantOk:     true,
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
			wantOk:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				gauge:   tt.storageState.gauge,
				counter: tt.storageState.counter,
			}
			got, ok, err := m.GetCounter(context.Background(), tt.metricName)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantOk, ok)
			assert.Equal(t, got, tt.want)
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
		wantOk       bool
	}{
		{
			name: "empty storage",
			storageState: storageState{
				gauge:   map[string]float64{},
				counter: map[string]int64{},
			},
			metricName: "testCounter",
			want:       0,
			wantOk:     false,
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
			wantOk:     false,
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
			wantOk:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				gauge:   tt.storageState.gauge,
				counter: tt.storageState.counter,
			}
			got, ok, err := m.GetGauge(context.Background(), tt.metricName)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantOk, ok)
			assert.Equal(t, tt.want, got)
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
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				gauge:   tt.storageState.gauge,
				counter: tt.storageState.counter,
			}
			m.IncCounter(ctx, tt.metric.name, tt.metric.value)
			actual, ok, err := m.GetCounter(ctx, tt.metric.name)
			assert.NoError(t, err)
			assert.True(t, ok)
			assert.Equal(t, tt.want, actual)
		})
	}
}

func TestMemStorage_SetGauge(t *testing.T) {
	m := NewMemStorage(nil, nil)
	assert.NotNil(t, m)
	ctx := context.Background()
	m.SetGauge(ctx, "test", 1)
	v, ok, err := m.GetGauge(ctx, "test")
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, 1.0, v)
}

func TestNewStorage(t *testing.T) {
	m := NewMemStorage(nil, nil)
	assert.NotNil(t, m)
}

func TestMemStorage_UpsertAll(t *testing.T) {
	one := int64(1)
	two := int64(2)
	pi := float64(3.14)
	twopi := float64(6.28)
	tests := []struct {
		name        string
		list        []models.Metrics
		wantGauge   map[string]float64
		wantCounter map[string]int64
	}{
		{
			name: "only counters",
			list: []models.Metrics{
				{ID: "asd", MType: models.MCounter, Delta: &two},
				{ID: "qwe", MType: models.MCounter, Delta: &one},
			},
			wantCounter: map[string]int64{"asd": two, "qwe": one},
			wantGauge:   map[string]float64{},
		},
		{
			name: "only float",
			list: []models.Metrics{
				{ID: "asd", MType: models.MGauge, Value: &pi},
				{ID: "qwe", MType: models.MGauge, Value: &twopi},
			},
			wantCounter: map[string]int64{},
			wantGauge:   map[string]float64{"asd": pi, "qwe": twopi},
		},
		{
			name: "counter and gauge",
			list: []models.Metrics{
				{ID: "asd", MType: models.MCounter, Delta: &one},
				{ID: "asd", MType: models.MGauge, Value: &pi},
				{ID: "qwe", MType: models.MCounter, Delta: &two},
				{ID: "qwe", MType: models.MGauge, Value: &twopi},
			},
			wantCounter: map[string]int64{"asd": one, "qwe": two},
			wantGauge:   map[string]float64{"asd": pi, "qwe": twopi},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMemStorage(nil, nil)
			require.NoError(t, m.UpsertAll(context.Background(), tt.list), fmt.Sprintf("UpsertAll(%v)", tt.list))
			require.Equal(t, tt.wantGauge, m.gauge)
			require.Equal(t, tt.wantCounter, m.counter)
		})
	}
}
