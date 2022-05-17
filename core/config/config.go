package config

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var defaultConfig = New()

func Default() *Config {
	return defaultConfig
}

func ResetDefault(c *Config) {
	defaultConfig = c

	Unmarshal = defaultConfig.Unmarshal
	Scan = defaultConfig.Scan
	Get = defaultConfig.Get
	GetString = defaultConfig.GetString
	GetBool = defaultConfig.GetBool
	GetInt = defaultConfig.GetInt
	GetFloat64 = defaultConfig.GetFloat64
	GetDuration = defaultConfig.GetDuration
	GetIntSlice = defaultConfig.GetIntSlice
	GetStringSlice = defaultConfig.GetStringSlice
	GetStringMap = defaultConfig.GetStringMap
}

type Config struct {
	v    *viper.Viper
	data []byte
	opts Options
}

func New(opts ...Option) *Config {
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

	callback := func() {
		for i := range c.opts.Callbacks {
			c.opts.Callbacks[i]()
		}
	}
	c.v.OnConfigChange(func(e fsnotify.Event) {
		data, _ := json.Marshal(c.v.AllSettings())
		if bytes.Compare(data, c.data) != 0 {
			c.data = data
			callback()
		}
	})
	c.v.WatchConfig()
	c.data, _ = json.Marshal(c.v.AllSettings())
	callback()
	return nil
}

func (c *Config) Unmarshal(val any) error {
	return c.v.Unmarshal(val)
}

func (c *Config) Scan(key string, val any) error {
	return c.v.UnmarshalKey(key, val)
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

var (
	Unmarshal      = defaultConfig.Unmarshal
	Scan           = defaultConfig.Scan
	Get            = defaultConfig.Get
	GetString      = defaultConfig.GetString
	GetBool        = defaultConfig.GetBool
	GetInt         = defaultConfig.GetInt
	GetFloat64     = defaultConfig.GetFloat64
	GetDuration    = defaultConfig.GetDuration
	GetIntSlice    = defaultConfig.GetIntSlice
	GetStringSlice = defaultConfig.GetStringSlice
	GetStringMap   = defaultConfig.GetStringMap
)
