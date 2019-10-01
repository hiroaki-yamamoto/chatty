package config

// Initializer of config

import (
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// New Initializes a configuration from the specified config name.
func New(cfgName string) (cfg *Config, err error) {
	var config Config
	viper.SetConfigName(cfgName)
	viper.AddConfigPath("/etc/real")
	viper.AddConfigPath("$HOME/etc/real")
	if err = viper.ReadInConfig(); err == nil {
		if err = viper.Unmarshal(&config); err == nil {
			cfg = &config
			cfg.Db.Client, err = mongo.NewClient(
				options.Client().ApplyURI(viper.GetString("db.URI")),
			)
			if err != nil {
				cfg = nil
				return
			}
			cfg.Db.Database = cfg.Db.Client.Database(viper.GetString("db.name"))
		}
	}
	return
}
