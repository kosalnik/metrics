package client

//import (
//	"testing"
//
//	"github.com/kosalnik/metrics/internal/config"
//	"github.com/kosalnik/metrics/internal/models"
//	"github.com/stretchr/testify/require"
//)
//
//func TestClient_collectMetrics(t *testing.T) {
//	c := NewClient(config.Agent{
//		Logger:           config.Logger{Level: "debug"},
//		CollectorAddress: "localhost",
//		PollInterval:     1,
//		ReportInterval:   1,
//	})
//	got := c.collectMetrics()
//	require.NotEmpty(t, got)
//	one := int64(1)
//	require.Contains(t, got, models.Metrics{ID: "PollCount", MType: models.MCounter, Delta: &one})
//
//}
