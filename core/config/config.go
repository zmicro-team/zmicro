package config

import (
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	DefaultConfig *config
)

type config struct {
	v    *viper.Viper
	opts Options
}

func NewConfig(opts ...Option) (*config, error) {
	options := Options{
		Type: "yaml",
		Path: "config.yaml",
	}

	for _, o := range opts {
		o(&options)
	}

	c := config{
		v:    viper.New(),
		opts: options,
	}

	c.Load()
	return &c, nil
}

func (c *config) Load() error {
	if c.opts.Type != "" {
		c.v.SetConfigType(c.opts.Type)
	}

	if c.opts.Path != "" {
		c.v.SetConfigFile(c.opts.Path)
	}

	if err := c.v.ReadInConfig(); err != nil {
		return err
	}

	c.v.OnConfigChange(func(e fsnotify.Event) {
	})
	c.v.WatchConfig()
	return nil
}

func (c *config) Scan(key string, val interface{}) error {
	return c.v.UnmarshalKey(key, val)
}

func (c *config) Get(key string) interface{} {
	return c.v.Get(key)
}

func (c *config) GetString(key string) string {
	return c.v.GetString(key)
}

func (c *config) GetBool(key string) bool {
	return c.v.GetBool(key)
}

func (c *config) GetInt(key string) int {
	return c.v.GetInt(key)
}

func (c *config) GetFloat64(key string) float64 {
	return c.v.GetFloat64(key)
}

func (c *config) GetDuration(key string) time.Duration {
	return c.v.GetDuration(key)
}

func (c *config) GetIntSlice(key string) []int {
	return c.v.GetIntSlice(key)
}

func (c *config) GetStringSlice(key string) []string {
	return c.v.GetStringSlice(key)
}

func (c *config) GetStringMap(key string) map[string]interface{} {
	return c.v.GetStringMap(key)
}

func Scan(key string, val interface{}) error {
	return DefaultConfig.Scan(key, val)
}

func Get(key string) interface{} {
	return DefaultConfig.Get(key)
}

func GetString(key string) string {
	return DefaultConfig.GetString(key)
}

func GetBool(key string) bool {
	return DefaultConfig.GetBool(key)
}

func GetInt(key string) int {
	return DefaultConfig.GetInt(key)
}

func GetFloat64(key string) float64 {
	return DefaultConfig.GetFloat64(key)
}

func GetDuration(key string) time.Duration {
	return DefaultConfig.GetDuration(key)
}

func GetIntSlice(key string) []int {
	return DefaultConfig.GetIntSlice(key)
}

func GetStringSlice(key string) []string {
	return DefaultConfig.GetStringSlice(key)
}

func GetStringMap(key string) map[string]interface{} {
	return DefaultConfig.GetStringMap(key)
}
