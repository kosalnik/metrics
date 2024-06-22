package memstorage

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
		want         *models.Metrics
	}{
		{
			name: "empty storage",
			storageState: storageState{
				gauge:   map[string]float64{},
				counter: map[string]int64{},
			},
			metricName: "testCounter",
			want:       nil,
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
			want:       &models.Metrics{ID: "testCounter", MType: models.MCounter, Delta: 3},
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
			want:       nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				gauge:   tt.storageState.gauge,
				counter: tt.storageState.counter,
			}
			got, err := m.GetCounter(context.Background(), tt.metricName)
			assert.NoError(t, err)
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
		want         *models.Metrics
	}{
		{
			name: "empty storage",
			storageState: storageState{
				gauge:   map[string]float64{},
				counter: map[string]int64{},
			},
			metricName: "testCounter",
			want:       nil,
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
			want:       nil,
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
			want:       &models.Metrics{ID: "testCounter", MType: models.MGauge, Value: 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				gauge:   tt.storageState.gauge,
				counter: tt.storageState.counter,
			}
			got, err := m.GetGauge(context.Background(), tt.metricName)
			assert.NoError(t, err)
			require.Equal(t, tt.want, got)
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
		want         *models.Metrics
	}{
		{
			name: "empty storage",
			storageState: storageState{
				gauge:   map[string]float64{},
				counter: map[string]int64{},
			},
			metric: metric{name: "test", value: 3},
			want:   &models.Metrics{ID: "test", MType: models.MCounter, Delta: 3},
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
			want:   &models.Metrics{ID: "test", MType: models.MCounter, Delta: 5},
		},
	}
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				gauge:   tt.storageState.gauge,
				counter: tt.storageState.counter,
			}
			r, err := m.IncCounter(ctx, tt.metric.name, tt.metric.value)
			require.NoError(t, err)
			require.Equal(t, tt.want, r)
			actual, err := m.GetCounter(ctx, tt.metric.name)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, actual)
		})
	}
}

func TestMemStorage_SetGauge(t *testing.T) {
	m := NewMemStorage()
	assert.NotNil(t, m)
	ctx := context.Background()
	want := &models.Metrics{ID: "test", MType: models.MGauge, Value: 1.0}
	r, err := m.SetGauge(ctx, "test", 1)
	require.NoError(t, err)
	require.Equal(t, want, r)
	v, err := m.GetGauge(ctx, "test")
	assert.NoError(t, err)
	assert.Equal(t, want, v)
}

func TestNewStorage(t *testing.T) {
	m := NewMemStorage()
	assert.NotNil(t, m)
}

func TestMemStorage_UpsertAll(t *testing.T) {
	tests := []struct {
		name        string
		list        []models.Metrics
		wantGauge   map[string]float64
		wantCounter map[string]int64
	}{
		{
			name: "only counters",
			list: []models.Metrics{
				{ID: "asd", MType: models.MCounter, Delta: 2},
				{ID: "qwe", MType: models.MCounter, Delta: 1},
			},
			wantCounter: map[string]int64{"asd": 2, "qwe": 1},
			wantGauge:   map[string]float64{},
		},
		{
			name: "only float",
			list: []models.Metrics{
				{ID: "asd", MType: models.MGauge, Value: 3.14},
				{ID: "qwe", MType: models.MGauge, Value: 6.28},
			},
			wantCounter: map[string]int64{},
			wantGauge:   map[string]float64{"asd": 3.14, "qwe": 6.28},
		},
		{
			name: "counter and gauge",
			list: []models.Metrics{
				{ID: "asd", MType: models.MCounter, Delta: 1},
				{ID: "asd", MType: models.MGauge, Value: 3.14},
				{ID: "qwe", MType: models.MCounter, Delta: 2},
				{ID: "qwe", MType: models.MGauge, Value: 6.28},
			},
			wantCounter: map[string]int64{"asd": 1, "qwe": 2},
			wantGauge:   map[string]float64{"asd": 3.14, "qwe": 6.28},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMemStorage()
			require.NoError(t, m.UpsertAll(context.Background(), tt.list), fmt.Sprintf("UpsertAll(%v)", tt.list))
			require.Equal(t, tt.wantGauge, m.gauge)
			require.Equal(t, tt.wantCounter, m.counter)
		})
	}
}

func BenchmarkMemStorage_UpsertAll(t *testing.B) {

	ctx := context.Background()
	m := NewMemStorage()
	for i := 0; i < t.N; i++ {
		require.NoError(t, m.UpsertAll(ctx, []models.Metrics{
			{ID: fmt.Sprintf("asd%d", i), MType: models.MCounter, Delta: int64(i)},
			{ID: fmt.Sprintf("asd%d", i), MType: models.MGauge, Value: 3.14},
			{ID: fmt.Sprintf("qwe%d", i), MType: models.MCounter, Delta: int64(i + 1)},
			{ID: fmt.Sprintf("qwe%d", i), MType: models.MGauge, Value: 6.28},
		}))
	}
}
