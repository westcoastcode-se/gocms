package event

type Bus struct {
	listeners []Listener
}

// Add a new event listener
func (b Bus) AddListener(listener Listener) {
	b.listeners = append(b.listeners, listener)
}

// Notify all listeners that a new event is sent. The event propagation will be aborted
// when the first listener returns an error
func (b *Bus) NotifyAll(e interface{}) error {
	for _, l := range b.listeners {
		if err := l.OnEvent(e); err != nil {
			return err
		}
	}
	return nil
}

// Create an ew event bus
func NewBus() *Bus {
	return &Bus{}
}
