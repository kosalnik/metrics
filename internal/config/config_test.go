package config

import (
	"reflect"
	"testing"

	"github.com/kosalnik/metrics/internal/infra/logger"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{name: "create", want: &Config{
			Agent: Agent{
				Logger:           logger.Config{Level: "info"},
				CollectorAddress: "127.0.0.1:8080",
				PollInterval:     2,
				ReportInterval:   10,
				RateLimit:        1,
			},
			Server: Server{
				Logger:  logger.Config{Level: "info"},
				Address: ":8080",
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
