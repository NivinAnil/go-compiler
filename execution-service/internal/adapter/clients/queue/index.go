package queue

type IQueueClient interface {
	Close() error
	ConsumeMessages(handleMessage func(string)) error
	PublishMessage(messageBody string) error
}
