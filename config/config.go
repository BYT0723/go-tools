package config

import (
	"context"
	"sync"

	"github.com/spf13/viper"
)

var (
	config *_config
	once   sync.Once
)

type _config struct {
	ctx            context.Context
	viper          *viper.Viper
	remote         bool
	decodeOpts     []viper.DecoderConfigOption
	onConfigChange ChangeHandler
}

func Init(opts ...Option) {
	once.Do(func() {
		config = &_config{
			ctx:    context.Background(),
			viper:  viper.New(),
			remote: false,
		}

		for _, opt := range opts {
			opt(config)
		}

		var err error
		if config.remote {
			err = initRemote()
		} else {
			err = initLocation()
		}

		if err != nil {
			panic(err)
		}

		if config.onConfigChange != nil {
			config.viper.OnConfigChange(config.onConfigChange)
		}
	})
}

func initLocation() (err error) {
	err = config.viper.ReadInConfig()
	if err != nil {
		return
	}
	config.viper.WatchConfig()
	return
}

func initRemote() (err error) {
	err = config.viper.ReadRemoteConfig()
	if err != nil {
		return
	}
	err = config.viper.WatchRemoteConfig()
	return
}

func Unmarshal(rawVal any) error {
	return config.viper.Unmarshal(rawVal, config.decodeOpts...)
}

func UnmarshalKey(key string, rawVal any) error {
	return config.viper.UnmarshalKey(key, rawVal, config.decodeOpts...)
}

func Viper() *viper.Viper {
	return config.viper
}
