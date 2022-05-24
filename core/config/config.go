package config

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var defaultConfig = New()

func Default() IConfig {
	return defaultConfig
}

func ResetDefault(c IConfig) {
	defaultConfig = c
}

type Config struct {
	v    *viper.Viper
	data []byte
	opts Options
}

func New(opts ...Option) IConfig {
	options := Options{
		Type: "yaml",
	}

	for _, o := range opts {
		o(&options)
	}

	c := Config{
		v:    viper.New(),
		opts: options,
	}

	if err := c.load(); err != nil {
		panic("load config error")
	}
	return &c
}

func (c *Config) load() error {
	if c.opts.Type != "" {
		c.v.SetConfigType(c.opts.Type)
	}

	if c.opts.Path == "" {
		return nil
	}

	c.v.SetConfigFile(c.opts.Path)

	if err := c.v.ReadInConfig(); err != nil {
		return err
	}

	callback := func(cfg IConfig) {
		for i := range c.opts.Callbacks {
			c.opts.Callbacks[i](cfg)
		}
	}
	c.v.OnConfigChange(func(e fsnotify.Event) {
		data, _ := json.Marshal(c.v.AllSettings())
		if bytes.Compare(data, c.data) != 0 {
			c.data = data
			callback(c)
		}
	})
	c.v.WatchConfig()
	c.data, _ = json.Marshal(c.v.AllSettings())
	callback(c)
	return nil
}

func (c *Config) Unmarshal(val any) error {
	temp, _ := json.Marshal(c.v.AllSettings())
	return json.Unmarshal(temp, val)
}

func (c *Config) Scan(key string, val any) error {
	temp, _ := json.Marshal(c.Get(key).(map[string]interface{}))
	return json.Unmarshal(temp, val)
}

func (c *Config) Get(key string) any {
	return c.v.Get(key)
}

func (c *Config) GetString(key string) string {
	return c.v.GetString(key)
}

func (c *Config) GetBool(key string) bool {
	return c.v.GetBool(key)
}

func (c *Config) GetInt(key string) int {
	return c.v.GetInt(key)
}

func (c *Config) GetFloat64(key string) float64 {
	return c.v.GetFloat64(key)
}

func (c *Config) GetDuration(key string) time.Duration {
	return c.v.GetDuration(key)
}

func (c *Config) GetIntSlice(key string) []int {
	return c.v.GetIntSlice(key)
}

func (c *Config) GetStringSlice(key string) []string {
	return c.v.GetStringSlice(key)
}

func (c *Config) GetStringMap(key string) map[string]any {
	return c.v.GetStringMap(key)
}

type IConfig interface {
	Unmarshal(val any) error
	Scan(key string, val any) error
	Get(key string) any
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetFloat64(key string) float64
	GetDuration(key string) time.Duration
	GetIntSlice(key string) []int
	GetStringSlice(key string) []string
	GetStringMap(key string) map[string]any
}
