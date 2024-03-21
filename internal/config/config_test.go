package config

import (
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{name: "create", want: &Config{
			Agent: Agent{
				Logger:           Logger{Level: "info"},
				CollectorAddress: "127.0.0.1:8080",
				PoolInterval:     2,
				ReportInterval:   10,
			},
			Server: Server{
				Logger:  Logger{Level: "info"},
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
