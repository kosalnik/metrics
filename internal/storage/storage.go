package storage

type Storage interface {
	GetGauge(name string) (float64, bool)
	SetGauge(name string, value float64)
	GetCounter(name string) (int64, bool)
	IncCounter(name string, value int64)
	GetPlain() map[string]string
}
