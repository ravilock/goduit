package config

import (
	"fmt"
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
	viper.SetDefault("log.level", "debug")
	viper.SetDefault("server.port", "3001")
	viper.SetDefault("server.domain", "localhost")
	viper.SetDefault("server.address", fmt.Sprintf("%s:%s", viper.GetString("server.domain"), viper.GetString("server.port")))
	viper.SetDefault("db.host", "mongodb")
	viper.SetDefault("db.port", "27017")
	viper.SetDefault("db.user", "goduit")
	viper.SetDefault("db.pass", "goduit-password")
	viper.SetDefault("db.url", "mongodb://goduit:goduit-password@mongo:27017/")
	viper.SetDefault("queue.type", "rabbitmq")
	viper.SetDefault("queue.url", "amqp://guest:guest@goduit-queue:5672/")
	viper.SetDefault("article.queue.name", "new-articles-queue")
	viper.SetDefault("feed.max.articles", 30)
}
