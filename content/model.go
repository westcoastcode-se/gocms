package content

import "time"

type Model struct {
	// Unique ID that represents this model
	ID string

	// Time when this instance was created
	CreatedAt time.Time

	// The view used when rendering this model
	View string

	// Type type
	Type string

	// The actual content
	Content interface{}
}
