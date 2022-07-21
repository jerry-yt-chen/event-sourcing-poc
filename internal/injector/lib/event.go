package lib

import (
	"github.com/jerry-yt-chen/event-sourcing-poc/configs"
	pubsub "github.com/jerry-yt-chen/event-sourcing-poc/pkg/event/pubsub"
	"github.com/jerry-yt-chen/event-sourcing-poc/pkg/fluentd"
)

func InitEventPublisher(mongoSvc fluentd.Service) pubsub.Publisher {
	return pubsub.NewPublisherDecorator(configs.C.Pub.ProjectID, configs.C.Pub.Topic, mongoSvc)
}
