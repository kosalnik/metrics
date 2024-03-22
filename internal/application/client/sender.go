package client

type Sender interface {
	SendGauge(k string, v float64)
	SendCounter(k string, v int64)
}
