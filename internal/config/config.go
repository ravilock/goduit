package config

import (
	"strings"

	"github.com/spf13/viper"
)

type customContextKey string

const (
	RequestIDContextKey     = customContextKey("X-Request-Id")
	AllowedGroupsContextKey = customContextKey("Allowed-Groups")
	IdentityContextKey      = customContextKey("Identity")
	FilterContextKey        = customContextKey("Filter")
)

func init() {
	viper.AutomaticEnv()
	envReplacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(envReplacer)
	viper.SetDefault("server.address", ":3000")
	viper.SetDefault("db.host", "mongodb")
	viper.SetDefault("db.port", "27017")
	viper.SetDefault("db.user", "goduit")
	viper.SetDefault("db.pass", "goduit-password")
	viper.SetDefault("db.url", "mongodb://goduit:goduit-password@mongo:27017/")
	viper.SetDefault("private.key.location", "/app/jwtRS256.key")
	viper.SetDefault("public.key.location", "/app/jwtRS256.key.pub")
}
