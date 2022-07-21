package persistence

import (
	"fmt"

	"github.com/jerry-yt-chen/event-sourcing-poc/configs"
	"github.com/jerry-yt-chen/event-sourcing-poc/pkg/fluentd"
)

func InitFluentd() fluentd.Service {
	client, err := fluentd.New(configs.C.Fluentd.EventLog)
	if err != nil {
		fmt.Printf("failed to connect to mongo err: %v\n", err)
		panic(err)
	}
	return client
}
