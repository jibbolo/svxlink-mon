package broker

// Loopback implements the Broker interface and it's used for testing
type Loopback struct {
	broker chan []byte
}

func NewLoopback() *Loopback {
	return &Loopback{
		make(chan []byte, 100),
	}
}

func (b *Loopback) Write(p []byte) (int, error) {
	b.broker <- p
	return 0, nil
}

func (b *Loopback) Close() error {
	close(b.broker)
	return nil
}

func (b *Loopback) Subscribe(fn MessageHandler) error {
	go func() {
		for msg := range b.broker {
			fn(msg)
		}
	}()
	return nil
}
