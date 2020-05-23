package cache

// Error for when a specific page is not found in the cache
type PageNotFound struct {
	Page string
}

func (p *PageNotFound) Error() string {
	return "could not find: " + p.Page
}
