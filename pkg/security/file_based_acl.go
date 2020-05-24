package security

import (
	"encoding/json"
	"github.com/westcoastcode-se/gocms/pkg/event"
	"io/ioutil"
	"log"
	"strings"
	"sync"
)

type aclEntry struct {
	Prefix string
	Roles  []string
}

type aclEntryBody struct {
	Access []aclEntry
}

type fileBasedACL struct {
	databasePath string
	mux          sync.Mutex
	Database     map[string]aclEntry
}

func (f *fileBasedACL) GetRoles(uri string) []string {
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

func (f *fileBasedACL) load() error {
	log.Printf(`Loading ACL from "%s"`+"\n", f.databasePath)
	bytes, err := ioutil.ReadFile(f.databasePath)
	if err != nil {
		return NewLoadError("Could not read database file: '%s' because: %e", f.databasePath, err)
	}

	var body aclEntryBody
	err = json.Unmarshal(bytes, &body)
	if err != nil {
		return NewLoadError("Could not parse database file: '%s' because: %e", f.databasePath, err)
	}

	var database = make(map[string]aclEntry)
	for _, a := range body.Access {
		database[a.Prefix] = a
	}

	f.mux.Lock()
	defer f.mux.Unlock()
	f.Database = database
	return nil
}

func (f *fileBasedACL) OnEvent(e interface{}) error {
	if _, ok := e.(*event.Checkout); ok {
		if err := f.load(); err != nil {
			return err
		}
	}
	return nil
}

// Create a new file-based ACL service.
func NewFileBasedACL(bus *event.Bus, path string) ACL {
	impl := &fileBasedACL{
		databasePath: path,
		mux:          sync.Mutex{},
		Database:     make(map[string]aclEntry),
	}
	if len(path) > 0 {
		err := impl.load()
		if err != nil {
			panic(err)
		}
	}
	bus.AddListener(impl)
	return impl
}
