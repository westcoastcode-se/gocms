package event

// Represents an event that a checkout has happened
type Checkout struct {
	// The commit that's been changed out
	Commit string
}

// Represents when changes are pushed to the remote server
type Push struct {
}
