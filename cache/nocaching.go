package cache

// Represents a cache that do not perform any caching
type NoCaching struct {
}

func (n NoCaching) Find(path string) ([]byte, error) {
	return nil, nil
}

func (n NoCaching) Set(path string, content []byte) {
}

func (n NoCaching) IsAllowed(path string) bool {
	return false
}

func (n NoCaching) Reset() {
}

func NewNoCaching() *NoCaching {
	return &NoCaching{}
}
