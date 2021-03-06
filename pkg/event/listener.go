package event

import "context"

type Listener interface {
	// This function will be called when someone is posting an event on the event bus. It's up to the
	// listener to capture the event if necessary.
	//
	// Return the underlying error if one occurs.
	OnEvent(ctx context.Context, e interface{}) error
}
