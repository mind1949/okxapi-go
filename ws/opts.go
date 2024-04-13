package ws

// OptFn client functional options
type OptFn func(*options)

func newDefaultOpts() options {
	return options{}
}

type options struct {
	messageHandler MessageHandler
}

// WithMessageHandler set message's handler
func WithMessageHandler(handler MessageHandler) OptFn {
	return func(o *options) {
		o.messageHandler = handler
	}
}
