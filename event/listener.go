package event

type Listener interface {
	OnEvent(e interface{}) error
}
