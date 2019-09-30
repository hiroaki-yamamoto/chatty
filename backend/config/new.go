package config

// Initializer of config

import "github.com/spf13/viper"

// New Initializes a configuration from the specified config name.
func New(cfgName string) (cfg *Config, err error) {
	var config Config
	viper.SetConfigName(cfgName)
	viper.AddConfigPath("/etc/real")
	viper.AddConfigPath("$HOME/etc/real")
	if err = viper.ReadInConfig(); err == nil {
		if err = viper.Unmarshal(&config); err == nil {
			cfg = &config
		}
	}
	return
}
