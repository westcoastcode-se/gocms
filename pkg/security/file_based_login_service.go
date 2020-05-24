package security

import (
	"encoding/base64"
	"encoding/json"
	"github.com/westcoastcode-se/gocms/pkg/event"
	"io/ioutil"
	"sync"
)

type userDatabaseUser struct {
	Username string
	Password string
	Roles    []string
}

type userDatabaseBody struct {
	Users []userDatabaseUser
}

type fileBasedLoginService struct {
	Users        []userDatabaseUser
	databasePath string
	mux          sync.Mutex

	publicKey  string
	privateKey string
}

func (s *fileBasedLoginService) OnEvent(e interface{}) error {
	if _, ok := e.(*event.Checkout); ok {
		if err := s.load(); err != nil {
			return err
		}
	}
	return nil
}

func (s *fileBasedLoginService) Login(username string, password string) (*User, error) {
	encoded := base64.StdEncoding.EncodeToString([]byte(password))
	s.mux.Lock()
	defer s.mux.Unlock()
	for _, user := range s.Users {
		if user.Username == username && user.Password == encoded {
			return &User{
				Name:  user.Username,
				Roles: user.Roles,
			}, nil
		}
	}
	return &User{
		Name:  "",
		Roles: []string{},
	}, &UserNotFound{username}
}

func (s *fileBasedLoginService) load() error {
	bytes, err := ioutil.ReadFile(s.databasePath)
	if err != nil {
		return NewLoadError("could not read database file: '%s' because: '%e'", s.databasePath, err)
	}

	var body userDatabaseBody
	err = json.Unmarshal(bytes, &body)
	if err != nil {
		return NewLoadError("could not parse database file: '%s' because: '%e'", s.databasePath, err)
	}

	var users []userDatabaseUser
	for _, u := range body.Users {
		users = append(users, userDatabaseUser{
			Username: u.Username,
			Password: u.Password,
			Roles:    u.Roles,
		})
	}

	s.mux.Lock()
	defer s.mux.Unlock()
	s.Users = users
	return nil
}

// Create a login service
func NewLoginService(bus *event.Bus, userDatabase string) LoginService {
	impl := &fileBasedLoginService{
		databasePath: userDatabase,
	}
	if len(userDatabase) > 0 {
		err := impl.load()
		if err != nil {
			panic(err)
		}
	}
	bus.AddListener(impl)
	return impl
}
