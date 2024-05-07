package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Option func(c *_config)

func WithConfigName(name string) Option {
	return func(c *_config) {
		c.viper.SetConfigName(name)
	}
}

func WithConfigType(t string) Option {
	return func(c *_config) {
		c.viper.SetConfigType(t)
	}
}

func WithConfigFile(file string) Option {
	return func(c *_config) {
		c.viper.SetConfigFile(file)
	}
}

func WithConfigPath(paths ...string) Option {
	return func(c *_config) {
		for _, path := range paths {
			c.viper.AddConfigPath(path)
		}
	}
}

func WithRemoteConfig(endpoint, path string) Option {
	return func(c *_config) {
		if err := c.viper.AddRemoteProvider("etcd", endpoint, path); err != nil {
			panic(err)
		}
		c.initStatus |= remote
	}
}

func OnConfigChange(h ChangeHandler) Option {
	return func(c *_config) {
		c.onConfigChange = h
	}
}

func WithConfigTag(name string) Option {
	return func(c *_config) {
		c.decodeOpts = append(c.decodeOpts, func(dc *mapstructure.DecoderConfig) {
			dc.TagName = name
		})
	}
}

func WithCustomDeocodeOpt(opts ...viper.DecoderConfigOption) Option {
	return func(c *_config) {
		c.decodeOpts = append(c.decodeOpts, opts...)
	}
}

func WithDefaultUnMarshal(unmarshaler Unmarshaler) Option {
	return func(c *_config) {
		c.unmarshaler = unmarshaler
	}
}
