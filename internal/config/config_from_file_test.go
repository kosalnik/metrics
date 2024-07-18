package config

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadFromJson(t *testing.T) {
	tests := []struct {
		name    string
		f       io.Reader
		c       Agent
		want    Agent
		wantErr bool
	}{
		{
			name: "empty json object",
			f:    strings.NewReader(`{}`),
			c:    Agent{CollectorAddress: "a", PollInterval: 99, ReportInterval: 999, RateLimit: 9999},
			want: Agent{CollectorAddress: "a", PollInterval: 99, ReportInterval: 999, RateLimit: 9999},
		},
		{
			name: "replace config with json defined",
			f:    strings.NewReader(`{"address": "newaddress:123"}`),
			c:    Agent{CollectorAddress: "a", PollInterval: 99, ReportInterval: 999, RateLimit: 9999},
			want: Agent{CollectorAddress: "newaddress:123", PollInterval: 99, ReportInterval: 999, RateLimit: 9999},
		},
		{
			name:    "wrong json",
			f:       strings.NewReader(`{`),
			c:       Agent{CollectorAddress: "a", PollInterval: 99, ReportInterval: 999, RateLimit: 9999},
			wantErr: true,
		},
		{
			name:    "wrong config",
			f:       strings.NewReader(`{"poll_interval":"one"}`),
			c:       Agent{CollectorAddress: "a", PollInterval: 99, ReportInterval: 999, RateLimit: 9999},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := loadFromJson(tt.f, &tt.c)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, tt.c)
			}
		})
	}
}
