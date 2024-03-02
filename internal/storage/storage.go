package storage

type Storage interface {
	GetGauge(name string) float64
	GetCounter(name string) int64
	SetGauge(name string, value float64)
	IncCounter(name string, value int64)
	GetPlain() map[string]string
}
