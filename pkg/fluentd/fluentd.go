package fluentd

type Service interface {
	Post(tag string, message interface{}) error
}
