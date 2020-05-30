package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/westcoastcode-se/gocms/pkg/event"
	"github.com/westcoastcode-se/gocms/pkg/security"
	"io/ioutil"
	"sync"
)

type userData struct {
	Username string
	Password string
	Roles    []string
}

type entryBody struct {
	Users []userData
}

type FileBasedLoginService struct {
	Users        []userData
	databasePath string
	mux          sync.Mutex

	publicKey  string
	privateKey string
}

func (s *FileBasedLoginService) OnEvent(_ context.Context, e interface{}) error {
	if _, ok := e.(*event.Checkout); ok {
		if err := s.load(); err != nil {
			return err
		}
	}
	return nil
}

func (s *FileBasedLoginService) Login(username string, password string) (*security.User, error) {
	encoded := base64.StdEncoding.EncodeToString([]byte(password))
	s.mux.Lock()
	defer s.mux.Unlock()
	for _, user := range s.Users {
		if user.Username == username && user.Password == encoded {
			return &security.User{
				Name:  user.Username,
				Roles: user.Roles,
			}, nil
		}
	}
	return &security.User{
		Name:  "",
		Roles: []string{},
	}, &UserNotFound{username}
}

func (s *FileBasedLoginService) load() error {
	bytes, err := ioutil.ReadFile(s.databasePath)
	if err != nil {
		return NewLoadError("could not read database file: '%s' because: '%e'", s.databasePath, err)
	}

	var body entryBody
	err = json.Unmarshal(bytes, &body)
	if err != nil {
		return NewLoadError("could not parse database file: '%s' because: '%e'", s.databasePath, err)
	}

	var users []userData
	for _, u := range body.Users {
		users = append(users, userData{
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
	impl := &FileBasedLoginService{
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
