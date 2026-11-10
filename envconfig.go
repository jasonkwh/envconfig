package envconfig

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

var _ EnvConfigurator = (*envConfiguratorImpl)(nil)

type envConfiguratorImpl struct {
	cfgFile *string
	cfg     any
}

func New(cfgFile *string, cfg any) EnvConfigurator {
	if cfg == nil {
		panic("cfg cannot be nil")
	}

	return &envConfiguratorImpl{
		cfgFile: cfgFile,
		cfg:     cfg,
	}
}

func (e *envConfiguratorImpl) Load() {
	if e.cfgFile != nil {
		viper.SetConfigFile(*e.cfgFile)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	bindEnvs(e.cfg, "")

	if err := viper.Unmarshal(e.cfg); err != nil {
		fmt.Printf("failed to unmarshal config: %v\n", err)
	}
}

func bindEnvs(iface any, prefix string) {
	if iface == nil {
		return
	}

	te := reflect.TypeOf(iface)

	if te.Kind() == reflect.Pointer {
		te = te.Elem()
	}

	// Only process struct types
	if te.Kind() != reflect.Struct {
		return
	}

	for field := range te.Fields() {
		tag := field.Tag.Get("mapstructure")
		if tag == "" {
			tag = strings.ToLower(field.Name)
		}

		if prefix != "" {
			tag = prefix + "." + tag
		}

		if field.Type.Kind() == reflect.Struct {
			bindEnvs(reflect.New(field.Type).Elem().Interface(), tag)
		} else {
			_ = viper.BindEnv(tag)
		}
	}
}
