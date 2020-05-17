package event

type Bus struct {
	listeners []Listener
}

func (b Bus) AddListener(listener Listener) {
	b.listeners = append(b.listeners, listener)
}

// Notify all listeners that a new event is sent
func (b *Bus) NotifyAll(e interface{}) error {
	for _, l := range b.listeners {
		if err := l.OnEvent(e); err != nil {
			return err
		}
	}
	return nil
}

func NewBus() *Bus {
	return &Bus{}
}
