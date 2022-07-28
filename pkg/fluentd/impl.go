package fluentd

import (
	"github.com/fluent/fluent-logger-golang/fluent"
)

type Config struct {
	Host string
	Port int
}

type impl struct {
	Client Service
}

func New(config Config) (Service, error) {
	client, err := fluent.New(fluent.Config{
		FluentHost: config.Host,
		FluentPort: config.Port,
		// If we want to make sure it must connect with Fluentd, please turn it to false.
		Async: true,
		// When async is enabled, this option defines the interval (ms) at which the connection to the fluentd-address is re-established.
		// This option is useful if the address may resolve to one or more IP addresses, e.g. a Consul service address.
		AsyncReconnectInterval: 60000,
	})
	if err != nil {
		panic(err)
	}

	return &impl{
		Client: client,
	}, nil
}

func (im *impl) Post(tag string, message interface{}) error {
	if err := im.Client.Post(tag, message); err != nil {
		return err
	}
	return nil
}
