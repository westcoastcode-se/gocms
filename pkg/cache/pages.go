package cache

type Pages interface {
	// Try to find a cached page. This will return an error if the cache do not contain the supplied page path
	Find(path string) ([]byte, error)

	// Set the cache for the supplied path
	Set(path string, content []byte)

	// Forcefully reset the cache
	Reset()

	// Check to see if the supplied path is allowed to be cached
	IsAllowed(path string) bool
}
