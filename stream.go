package jsonstream

type (
	Receiver interface {
		Receive(data any) error
	}

	Sender interface {
		Send(data any) error
		Done() error
	}
)
