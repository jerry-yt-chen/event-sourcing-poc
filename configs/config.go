package configs

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/jerry-yt-chen/event-sourcing-poc/pkg/mongo"
)

const (
	fileType = "yaml"
)

var (
	C = new(Config)
)

type Config struct {
	App   AppConfig
	Pub   PubSubConfig `mapstructure:"publisher"`
	Sub   PubSubConfig `mapstructure:"subscriber"`
	Mongo mongo.Config
}

type AppConfig struct {
	BaseRoute     string
	Port          int
	RunMode       string
	Name          string
	EnableProfile bool
}

type PubSubConfig struct {
	ProjectID string
	Topic     string
}

func InitConfigs() {
	fromFile()
	fmt.Printf("C.App = %+v\n", C.App)
	fmt.Printf("C.Pub = %+v\n", C.Pub)
	fmt.Printf("C.Sub = %+v\n", C.Sub)
	fmt.Printf("C.Mongo = %+v\n", C.Mongo)
}

func fromFile() {
	viper.SetConfigType(fileType)
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Read configs error，reason：" + err.Error())
	}

	err = viper.Unmarshal(&C)
	if err != nil {
		panic("Unmarshal error: " + err.Error())
	}
}
