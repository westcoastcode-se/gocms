package content

type Navigation struct {
	URI string
}

func (n *Navigation) Active(uri string) string {
	if uri == n.URI {
		return "active"
	}
	return ""
}

// Check too see if the supplied uri matches the current uri. Will return true
// if this is the case.
func (n *Navigation) Matches(uri string) bool {
	return uri == n.URI
}
