package content

// Service that's responsible for all the content managed by the cms
type Controller interface {
	// Update the content managed by this controller
	Update(commit string) error

	// Save the content managed by this controller
	Save(message string) error
}
