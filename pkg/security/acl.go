package security

type ACL interface {
	// Fetch roles required for accessing the supplied uri.
	GetRoles(uri string) []string
}
