package cfg

import (
	"sync"

	"github.com/spf13/viper"
)

var (
	config *_config
	once   sync.Once
)

type (
	initStatus  uint8
	Unmarshaler func(*viper.Viper) error
)

const (
	remote initStatus = 1 << iota
)

type _config struct {
	viper          *viper.Viper
	decodeOpts     []viper.DecoderConfigOption
	initStatus     initStatus
	onConfigChange ChangeHandler
	unmarshaler    Unmarshaler
}

func Init(opts ...Option) {
	once.Do(func() {
		config = &_config{
			viper: viper.New(),
		}

		for _, opt := range opts {
			opt(config)
		}

		defer func() {
		}()

		var err error
		if config.initStatus&remote == remote {
			if err = config.viper.ReadRemoteConfig(); err != nil {
				panic(err)
			}

			if config.onConfigChange != nil {
				if err = config.viper.WatchRemoteConfig(); err != nil {
					panic(err)
				}
				config.viper.OnConfigChange(config.onConfigChange)
			}
		} else {
			if err = config.viper.ReadInConfig(); err != nil {
				panic(err)
			}
			if config.onConfigChange != nil {
				config.viper.WatchConfig()
				config.viper.OnConfigChange(config.onConfigChange)
			}
		}

		if config.unmarshaler != nil {
			if err := config.unmarshaler(config.viper); err != nil {
				panic(err)
			}
		}
	})
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
