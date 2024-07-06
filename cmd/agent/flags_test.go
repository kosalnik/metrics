package main

import (
	"os"
	"testing"

	"github.com/kosalnik/metrics/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseFlags(t *testing.T) {
	cases := []struct {
		name  string
		flags []string
		env   map[string]string
		want  func(c *config.Agent)
	}{
		{
			name:  "CollectorAddress in flag",
			flags: []string{"script", "-a=example.com:123"},
			want: func(c *config.Agent) {
				assert.Equal(t, "example.com:123", c.CollectorAddress)
			},
		},
		{
			name:  "CollectorAddress in ENV",
			flags: []string{"script"},
			env:   map[string]string{"ADDRESS": "example.com:123"},
			want: func(c *config.Agent) {
				assert.Equal(t, "example.com:123", c.CollectorAddress)
			},
		},
		{
			name:  "CollectorAddress in flags and ENV",
			flags: []string{"script", "-a=example.com:123"},
			env:   map[string]string{"ADDRESS": "example.com:456"},
			want: func(c *config.Agent) {
				assert.Equal(t, "example.com:456", c.CollectorAddress)
			},
		},
		{
			name:  "REPORT_INTERVAL in flag",
			flags: []string{"script", "-r=33"},
			want: func(c *config.Agent) {
				assert.Equal(t, int64(33), c.ReportInterval)
			},
		},
		{
			name:  "REPORT_INTERVAL in ENV",
			flags: []string{"script"},
			env:   map[string]string{"REPORT_INTERVAL": "11"},
			want: func(c *config.Agent) {
				assert.Equal(t, int64(11), c.ReportInterval)
			},
		},
		{
			name:  "REPORT_INTERVAL in flags and ENV",
			flags: []string{"script", "-r=33"},
			env:   map[string]string{"REPORT_INTERVAL": "11"},
			want: func(c *config.Agent) {
				assert.Equal(t, int64(11), c.ReportInterval)
			},
		},
		{
			name:  "RATE_LIMIT in flag",
			flags: []string{"script", "-l=13"},
			want: func(c *config.Agent) {
				assert.Equal(t, int64(13), c.RateLimit)
			},
		},
		{
			name:  "RATE_LIMIT in ENV",
			flags: []string{"script"},
			env:   map[string]string{"RATE_LIMIT": "34"},
			want: func(c *config.Agent) {
				assert.Equal(t, int64(34), c.RateLimit)
			},
		},
		{
			name:  "RATE_LIMIT in flags and ENV",
			flags: []string{"script", "-l=13"},
			env:   map[string]string{"RATE_LIMIT": "55"},
			want: func(c *config.Agent) {
				assert.Equal(t, int64(55), c.RateLimit)
			},
		},
		{
			name:  "PROFILING in ENV",
			flags: []string{"script"},
			env:   map[string]string{"PROFILING": "true"},
			want: func(c *config.Agent) {
				assert.True(t, c.Profiling.Enabled)
			},
		},
		{
			name:  "POLL_INTERVAL in ENV",
			flags: []string{"script"},
			env:   map[string]string{"POLL_INTERVAL": "123"},
			want: func(c *config.Agent) {
				assert.Equal(t, int64(123), c.PollInterval)
			},
		},
		{
			name:  "KEY false in ENV",
			flags: []string{"script"},
			env:   map[string]string{"KEY": "jasdjhfqwehriuh"},
			want: func(c *config.Agent) {
				assert.Equal(t, "jasdjhfqwehriuh", c.Hash.Key)
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			c := config.NewAgent()
			if tt.env == nil {
				parseFlags(tt.flags, c)
			} else {
				old := make(map[string]string, len(tt.env))
				for k, v := range tt.env {
					old[k] = os.Getenv(k)
					require.NoError(t, os.Setenv(k, v))
				}
				parseFlags(tt.flags, c)
				for k, v := range old {
					require.NoError(t, os.Setenv(k, v))
				}
			}
			if tt.want != nil {
				tt.want(c)
			}
		})
	}
}
