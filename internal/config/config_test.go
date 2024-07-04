package config

import (
	"reflect"
	"testing"

	"github.com/kosalnik/metrics/internal/logger"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		want *Config
		name string
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

func TestNewServer(t *testing.T) {
	tests := []struct {
		want *Server
		name string
	}{
		{name: "create", want: &Server{
			Logger:  logger.Config{Level: "info"},
			Address: ":8080",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewServer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAgent(t *testing.T) {
	tests := []struct {
		want *Agent
		name string
	}{
		{name: "create", want: &Agent{
			Logger:           logger.Config{Level: "info"},
			CollectorAddress: "127.0.0.1:8080",
			PollInterval:     2,
			ReportInterval:   10,
			RateLimit:        1,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAgent(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAgent() = %v, want %v", got, tt.want)
			}
		})
	}
}
