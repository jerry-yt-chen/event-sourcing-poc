package event

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/sirupsen/logrus"

	"github.com/jerry-yt-chen/event-sourcing-poc/pkg/event"
	"github.com/jerry-yt-chen/event-sourcing-poc/pkg/fluentd"
)

type PublisherDecorator struct {
	pub     Publisher
	fluentd fluentd.Service
	topic   string
	tag     string
}

func NewPublisherDecorator(projectID, topic string, fluentd fluentd.Service) Publisher {
	publisher, err := NewGcpPublisher(projectID, topic)
	if err != nil {
		panic(err)
	}
	return &PublisherDecorator{
		pub:     publisher,
		fluentd: fluentd,
		topic:   topic,
		tag:     fmt.Sprintf("event.%s.publisher", topic),
	}
}

func (d *PublisherDecorator) Send(payload interface{}, metadata event.Metadata) error {
	var eventID string
	if id, ok := metadata["eventID"]; ok {
		eventID = id
	} else {
		eventID = watermill.NewUUID()
		metadata["eventID"] = eventID
	}

	if err := d.pub.Send(payload, metadata); err != nil {
		return err
	}
	d.saveRecord(payload, metadata)
	return nil
}

func (d *PublisherDecorator) saveRecord(payload interface{}, metadata event.Metadata) {
	p, _ := json.Marshal(payload)
	record := event.PublishRecord{
		TraceID:     metadata[event.TraceAttribute],
		Topic:       d.topic,
		EventID:     metadata["eventID"],
		Payload:     string(p),
		PublishTime: time.Now().Unix(),
		CreatedTime: time.Now().Unix(),
	}
	if err := d.fluentd.Post(d.tag, record); err != nil {
		logrus.WithField("err", err).Error("Insert record failed")
	}
}
