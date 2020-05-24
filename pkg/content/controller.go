package content

import "context"

// Service that's responsible for all the content managed by the cms
type Controller interface {
	// Update the content managed by this controller
	Update(ctx context.Context, commit string) error

	// Save the content managed by this controller
	Save(ctx context.Context, message string) error
}
