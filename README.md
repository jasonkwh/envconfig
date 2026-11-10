# envconfig

<a href="https://github.com/jasonkwh/envconfig/actions/workflows/coverage.yml"><img src="https://img.shields.io/github/actions/workflow/status/jasonkwh/envconfig/coverage.yml?label=coverage" alt="coverage"></a>
<a href="https://github.com/jasonkwh/envconfig/actions/workflows/lint.yml"><img src="https://img.shields.io/github/actions/workflow/status/jasonkwh/envconfig/lint.yml?label=lint" alt="lint"></a>
<a href="https://github.com/jasonkwh/envconfig/actions/workflows/test.yml"><img src="https://img.shields.io/github/actions/workflow/status/jasonkwh/envconfig/test.yml?label=test" alt="test"></a>
<a href="https://github.com/jasonkwh/envconfig/actions/workflows/vet.yml"><img src="https://img.shields.io/github/actions/workflow/status/jasonkwh/envconfig/vet.yml?label=vet" alt="vet"></a>
<a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="MIT License"></a>

`envconfig` loads `config.yaml` into a struct. Environment variables for the same fields take priority over the file.

This is an example of initialisation:

```go
type Config struct {
	Port        int
	MinIdleTime time.Duration
	Server      struct {
		Port int
	}
}

var (
	cfgFile string
	cfg     Config
)

func init() {
	cfgr := envconfig.New(&cfgFile, &cfg)
	cobra.OnInitialize(cfgr.Load)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.yaml", "config file")
}
```

`Load` runs after flags are parsed and maps the file into `cfg`. A field name is lowercased, nested fields are joined with `.`, and the environment variable is that key uppercased with `.` replaced by `_`. `Server.Port` reads `SERVER_PORT`.

```bash
go get github.com/jasonkwh/envconfig
```
