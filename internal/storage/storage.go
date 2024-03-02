package storage

type Storage interface {
	HasGauge(name string) bool
	GetGauge(name string) float64
	SetGauge(name string, value float64)
	HasCounter(name string) bool
	GetCounter(name string) int64
	IncCounter(name string, value int64)
	GetPlain() map[string]string
}
