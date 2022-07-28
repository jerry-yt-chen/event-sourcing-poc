package main

import (
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/sirupsen/logrus"

	"github.com/jerry-yt-chen/event-sourcing-poc/configs"
	"github.com/jerry-yt-chen/event-sourcing-poc/pkg/event"
	pubsub "github.com/jerry-yt-chen/event-sourcing-poc/pkg/event/pubsub"
	"github.com/jerry-yt-chen/event-sourcing-poc/pkg/fluentd"
)

var tag string

func main() {
	configs.InitConfigs()
	fluentdSvc, _ := fluentd.New(configs.C.Fluentd.EventLog)
	subscriber, _ := pubsub.NewGcpSubscriber(configs.C.Sub)

	// Subscribe will create the subscription. Only messages that are sent after the subscription is created may be received.
	messages, err := subscriber.Subscribe(configs.C.Sub.Topic)
	if err != nil {
		panic(err)
	}

	tag = fmt.Sprintf("event.%s.subscriber", configs.C.Sub.Topic)
	process(messages, fluentdSvc)
}

func process(messages <-chan *message.Message, mongoSvc fluentd.Service) {
	for msg := range messages {
		logrus.Printf("received id: %s, event: %s, publishedTime: %v\n", msg.UUID, string(msg.Payload), msg.Metadata.Get("publishTime"))
		receiveTime := time.Now()
		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
		go saveRecord(msg, mongoSvc, receiveTime)
	}
}

func saveRecord(m *message.Message, fluentdSvc fluentd.Service, receiveTime time.Time) {
	traceID := m.Metadata.Get("Cloud-Trace-Context")
	record := event.ReceiveRecord{
		Topic:       configs.C.Sub.Topic,
		TraceID:     traceID,
		EventID:     m.UUID,
		ReceiveTime: receiveTime.Unix(),
		CreatedTime: time.Now().Unix(),
	}
	if err := fluentdSvc.Post(tag, record); err != nil {
		logrus.WithField("err", err).Error("Post failed")
	}
}
