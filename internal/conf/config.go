package conf

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Server     ServerConfig
	Datasource DataSourceConfig
	Image      ImgConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type DataSourceConfig struct {
	Postgres PostgresConfig
	Redis    RedisConfig
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type RedisConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database int
}

type ImgConfig struct {
	SvrUrl    string
	DirectUrl string
	Auth      string
}

var Conf *Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("read Conf failed: %v", err)
	}

	Conf = &Config{}
	err = viper.Unmarshal(Conf)
	if err != nil {
		log.Fatalf("unmarshal Conf failed: %v", err)
	}

	log.Printf("Conf: %+v", Conf)

	// 监听配置变化
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	viper.WatchConfig()
}

func GetServerAddr() string {
	return Conf.Server.Host + ":" + fmt.Sprintf("%d", Conf.Server.Port)
}

func GetDBConfig() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		Conf.Datasource.Postgres.Host,
		Conf.Datasource.Postgres.User,
		Conf.Datasource.Postgres.Password,
		Conf.Datasource.Postgres.Database,
		Conf.Datasource.Postgres.Port)
}

func GetRedisConfig() string {
	//	"redis://<user>:<pass>@localhost:6379/<db>"
	return fmt.Sprintf("redis://%s:%s@%s:%d/%d",
		Conf.Datasource.Redis.User,
		Conf.Datasource.Redis.Password,
		Conf.Datasource.Redis.Host,
		Conf.Datasource.Redis.Port,
		Conf.Datasource.Redis.Database)

}

func GetImgConfig() string {
	return fmt.Sprintf("svrUrl: %s; directUrl: %s; auth string: %s;",
		Conf.Image.SvrUrl,
		Conf.Image.DirectUrl,
		Conf.Image.Auth)
}
