package acl

import (
	"context"
	"encoding/json"
	"github.com/westcoastcode-se/gocms/pkg/event"
	"github.com/westcoastcode-se/gocms/pkg/log"
	"io/ioutil"
	"strings"
	"sync"
)

type entry struct {
	Prefix string
	Roles  []string
}

type entryBody struct {
	Access []entry
}

type DefaultService struct {
	databasePath string
	mux          sync.Mutex
	Database     map[string]entry
}

func (f *DefaultService) GetRoles(uri string) []string {
	var distance = 100000
	var roles []string

	f.mux.Lock()
	defer f.mux.Unlock()
	for key, value := range f.Database {
		if len(key) > len(uri) {
			continue
		}

		if strings.HasPrefix(uri, key) {
			var newDistance = len(uri) - len(key)
			if newDistance < distance {
				distance = newDistance
				roles = value.Roles
			}
		}
	}
	return roles
}

func (f *DefaultService) load(ctx context.Context) error {
	log.Infof(ctx, "Loading ACL from %s", f.databasePath)
	bytes, err := ioutil.ReadFile(f.databasePath)
	if err != nil {
		return NewLoadError("Could not read database file: '%s' because: %e", f.databasePath, err)
	}

	var body entryBody
	err = json.Unmarshal(bytes, &body)
	if err != nil {
		return NewLoadError("Could not parse database file: '%s' because: %e", f.databasePath, err)
	}

	var database = make(map[string]entry)
	for _, a := range body.Access {
		database[a.Prefix] = a
	}

	f.mux.Lock()
	defer f.mux.Unlock()
	f.Database = database
	return nil
}

func (f *DefaultService) OnEvent(ctx context.Context, e interface{}) error {
	if _, ok := e.(*event.Checkout); ok {
		if err := f.load(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Create a new file-based ACL service.
func NewFileBasedACL(bus *event.Bus, path string) Service {
	impl := &DefaultService{
		databasePath: path,
		mux:          sync.Mutex{},
		Database:     make(map[string]entry),
	}
	if len(path) > 0 {
		err := impl.load(context.Background())
		if err != nil {
			panic(err)
		}
	}
	bus.AddListener(impl)
	return impl
}
