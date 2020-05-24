package cache

import (
	"context"
	"encoding/json"
	"github.com/westcoastcode-se/gocms/pkg/event"
	"io/ioutil"
	"log"
	"strings"
	"sync"
)

type cacheDatabaseBody struct {
	Whitelist []string
	Blacklist []string
}

// Represents a cache that saves page content indefinitely until Reset is called.
// Reset is normally called when a new sync request happens
type PermanentCache struct {
	mux          sync.Mutex
	data         map[string][]byte
	databasePath string
	database     cacheDatabaseBody
}

func (p *PermanentCache) Find(path string) ([]byte, error) {
	p.mux.Lock()
	defer p.mux.Unlock()
	if c, ok := p.data[path]; ok {
		return c, nil
	}
	return []byte{}, &PageNotFound{path}
}

func (p *PermanentCache) Set(path string, content []byte) {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.data[path] = content
}

func (p *PermanentCache) Reset() {
	log.Println("Resetting cache")
	p.mux.Lock()
	defer p.mux.Unlock()
	p.data = make(map[string][]byte)
}

func (p *PermanentCache) OnEvent(ctx context.Context, e interface{}) error {
	if _, ok := e.(*event.Checkout); ok {
		err := p.load()
		if err != nil {
			log.Printf("Could not load page cache database. Reason: %e\n", err)
			return err
		}
		p.Reset()
	}
	return nil
}

func (p *PermanentCache) IsAllowed(path string) bool {
	p.mux.Lock()
	defer p.mux.Unlock()

	var whitelisted = false
	for _, prefix := range p.database.Whitelist {
		if strings.HasPrefix(path, prefix) {
			whitelisted = true
			break
		}
	}

	if !whitelisted {
		return false
	}

	for _, prefix := range p.database.Blacklist {
		if strings.HasPrefix(path, prefix) {
			return false
		}
	}

	return true
}

func (p *PermanentCache) load() error {
	log.Printf(`Loading cache database from "%s"`+"\n", p.databasePath)
	bytes, err := ioutil.ReadFile(p.databasePath)
	if err != nil {
		log.Printf(`Could not read database file "%s". Reason: %e\n`, p.databasePath, err)
		return err
	}

	var body cacheDatabaseBody
	err = json.Unmarshal(bytes, &body)
	if err != nil {
		log.Printf(`Could not parse database "%s". Reason: %e\n`, p.databasePath, err)
		return err
	}

	p.mux.Lock()
	defer p.mux.Unlock()
	p.database = body
	return nil
}

// Create a new permanent cache. The only way to reset the cache is when the underlying content
// is updated. This is managed by listening for specific events on the event bus.
func NewPermanentCache(bus *event.Bus, databasePath string) *PermanentCache {
	impl := &PermanentCache{
		data:         make(map[string][]byte),
		databasePath: databasePath,
	}
	if len(databasePath) > 0 {
		err := impl.load()
		if err != nil {
			panic(err)
		}
	}
	bus.AddListener(impl)
	return impl
}
