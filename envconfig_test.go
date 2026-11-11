package envconfig_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/jasonkwh/envconfig"
	"github.com/spf13/viper"
)

type loadConfig struct {
	Port        int
	MinIdleTime time.Duration
	Server      struct {
		Port        int
		MinIdleTime time.Duration
	}
}

func TestEnvConfigurator_Load(t *testing.T) {
	tests := []struct {
		name     string
		cfgFile  string
		env      map[string]string
		cfg      loadConfig
		expected loadConfig
	}{
		{
			name:    "default",
			cfgFile: writeConfigFile(t, "Port: 8080"),
			expected: loadConfig{
				Port: 8080,
			},
		},
		{
			name:    "reading time from config file",
			cfgFile: writeConfigFile(t, "Port: 8080\nMinIdleTime: 10s"),
			expected: loadConfig{
				Port:        8080,
				MinIdleTime: 10 * time.Second,
			},
		},
		{
			name:    "reading nested struct from config file",
			cfgFile: writeConfigFile(t, "Server:\n  Port: 8080\n  MinIdleTime: 10s"),
			expected: loadConfig{
				Server: struct {
					Port        int
					MinIdleTime time.Duration
				}{
					Port:        8080,
					MinIdleTime: 10 * time.Second,
				},
			},
		},
		{
			name: "reading from environment variables",
			env:  map[string]string{"PORT": "8080"},
			expected: loadConfig{
				Port: 8080,
			},
		},
		{
			name: "reading nested struct from environment variables",
			env:  map[string]string{"SERVER_PORT": "80", "SERVER_MINIDLETIME": "5s"},
			expected: loadConfig{
				Server: struct {
					Port        int
					MinIdleTime time.Duration
				}{
					Port:        80,
					MinIdleTime: 5 * time.Second,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			viper.Reset()

			if test.env != nil {
				for key, value := range test.env {
					t.Setenv(key, value)
				}
			}

			as := assert.New(t)
			cfg := test.cfg
			cfgr := envconfig.New(&test.cfgFile, &cfg)

			cfgr.Load()
			as.Equal(test.expected, cfg)
		})
	}
}

func writeConfigFile(t *testing.T, content string) string {
	t.Helper()

	cfgFile := filepath.Join(t.TempDir(), "config.yaml")

	if err := os.WriteFile(cfgFile, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	return cfgFile
}
