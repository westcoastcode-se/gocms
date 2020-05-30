package acl

// Service used to fetch which roles are required to access specific URL's
type Service interface {
	// Fetch roles required for accessing the supplied uri.
	GetRoles(uri string) []string
}
