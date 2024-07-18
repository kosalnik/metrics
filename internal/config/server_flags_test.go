package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseServerFlags(t *testing.T) {
	cases := []struct {
		name    string
		flags   []string
		env     map[string]string
		want    func(c *Server)
		wantErr bool
	}{
		{
			name:  "Address in flag",
			flags: []string{"script", "-a=example.com:123"},
			want: func(c *Server) {
				assert.Equal(t, "example.com:123", c.Address)
			},
		},
		{
			name:  "Address in ENV",
			flags: []string{"script"},
			env:   map[string]string{"ADDRESS": "example.com:123"},
			want: func(c *Server) {
				assert.Equal(t, "example.com:123", c.Address)
			},
		},
		{
			name:  "Address in flags and ENV",
			flags: []string{"script", "-a=example.com:123"},
			env:   map[string]string{"ADDRESS": "example.com:456"},
			want: func(c *Server) {
				assert.Equal(t, "example.com:456", c.Address)
			},
		},
		{
			name:  "DATABASE_DSN in flag",
			flags: []string{"script", "-d=example.com:123"},
			want: func(c *Server) {
				assert.Equal(t, "example.com:123", c.DB.DSN)
			},
		},
		{
			name:  "DATABASE_DSN in ENV",
			flags: []string{"script"},
			env:   map[string]string{"DATABASE_DSN": "example.com:123"},
			want: func(c *Server) {
				assert.Equal(t, "example.com:123", c.DB.DSN)
			},
		},
		{
			name:  "DATABASE_DSN in flags and ENV",
			flags: []string{"script", "-d=example.com:123"},
			env:   map[string]string{"DATABASE_DSN": "example.com:456"},
			want: func(c *Server) {
				assert.Equal(t, "example.com:456", c.DB.DSN)
			},
		},
		{
			name:  "StoreInterval in flag",
			flags: []string{"script", "-i=13"},
			want: func(c *Server) {
				assert.Equal(t, 13, c.Backup.StoreInterval)
			},
		},
		{
			name:  "StoreInterval in ENV",
			flags: []string{"script"},
			env:   map[string]string{"STORE_INTERVAL": "34"},
			want: func(c *Server) {
				assert.Equal(t, 34, c.Backup.StoreInterval)
			},
		},
		{
			name:  "StoreInterval in flags and ENV",
			flags: []string{"script", "-i=13"},
			env:   map[string]string{"STORE_INTERVAL": "55"},
			want: func(c *Server) {
				assert.Equal(t, 55, c.Backup.StoreInterval)
			},
		},
		{
			name:  "PROFILING in ENV",
			flags: []string{"script"},
			env:   map[string]string{"PROFILING": "true"},
			want: func(c *Server) {
				assert.True(t, c.Profiling.Enabled)
			},
		},
		{
			name:  "FILE_STORAGE_PATH in ENV",
			flags: []string{"script"},
			env:   map[string]string{"FILE_STORAGE_PATH": "/z/x/c"},
			want: func(c *Server) {
				assert.Equal(t, "/z/x/c", c.Backup.FileStoragePath)
			},
		},
		{
			name:  "RESTORE true in ENV",
			flags: []string{"script"},
			env:   map[string]string{"RESTORE": "true"},
			want: func(c *Server) {
				assert.True(t, c.Backup.Restore)
			},
		},
		{
			name:  "RESTORE false in ENV",
			flags: []string{"script"},
			env:   map[string]string{"RESTORE": "false"},
			want: func(c *Server) {
				assert.False(t, c.Backup.Restore)
			},
		},
		{
			name:  "KEY false in ENV",
			flags: []string{"script"},
			env:   map[string]string{"KEY": "jasdjhfqwehriuh"},
			want: func(c *Server) {
				assert.Equal(t, "jasdjhfqwehriuh", c.Hash.Key)
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			c := NewServer()
			if tt.env == nil {
				require.NoError(t, ParseServerFlags(tt.flags, c))
			} else {
				old := make(map[string]string, len(tt.env))
				for k, v := range tt.env {
					old[k] = os.Getenv(k)
					require.NoError(t, os.Setenv(k, v))
				}
				require.NoError(t, ParseServerFlags(tt.flags, c))
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
