package content

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/westcoastcode-se/gocms/pkg/event"
	"github.com/westcoastcode-se/gocms/pkg/log"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type pageData struct {
	ID        string
	CreatedAt time.Time
	View      string
	Type      string
	Content   json.RawMessage
}

// Function for unmarshal the actual content
type UnmarshalContentFunc func(msg json.RawMessage) (interface{}, error)

//
type Repository interface {
	// Register a new type of model. For example:
	//  repository.RegisterModelType("models.News", models.JsonToNews)
	RegisterModelType(name string, fn UnmarshalContentFunc)

	// Save the supplied model. If the save failed for some reason then an error will be returned
	Save(ctx context.Context, path string, model *Model) (*Model, error)

	// Search for the model associated with the supplied path. The repository will return an error and the default
	// 404 model if the supplied path is not found.
	FindByPath(path string) (*Model, error)

	//
	Reload(ctx context.Context) error

	// Search for content with the supplied content type.
	Search(contentType string) []*SearchResult

	// Search for content with the supplied uuid
	Lookup(uuid string) *SearchResult

	// Fetch all
	GetAll() []*SearchResult
}

type RepositoryImpl struct {
	rootPath string
	mux      sync.Mutex
	Data     map[string]*Model
	Types    map[string]UnmarshalContentFunc
}

func (r *RepositoryImpl) RegisterModelType(view string, fn UnmarshalContentFunc) {
	r.Types[view] = fn
}

func (r *RepositoryImpl) Save(ctx context.Context, p string, model *Model) (*Model, error) {
	var id = model.ID
	if id == "" {
		id = uuid.New().String()
	}

	var contentJson json.RawMessage = nil
	contentJson, err := json.Marshal(model.Content)
	if err != nil {
		return nil, err
	}

	output := pageData{
		id,
		model.CreatedAt,
		model.View,
		model.Type,
		contentJson,
	}
	b, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}

	absolutePath := path.Join(r.rootPath, p)
	err = ioutil.WriteFile(absolutePath+".json", b, 0644)
	if err != nil {
		return nil, err
	}

	log.Infof(ctx, "Sucessfully saved %s", p)
	model.ID = id
	return model, nil
}

func (r *RepositoryImpl) FindByPath(path string) (*Model, error) {
	r.mux.Lock()
	defer r.mux.Unlock()
	if val, found := r.Data[path]; found {
		return val, nil
	}

	return &Model{
		View:    "views/errors/404.html",
		Content: nil,
	}, NewNotFoundError(path)
}

func (r *RepositoryImpl) Reload(ctx context.Context) error {
	log.Infof(ctx, "Reloading content from dir %s", r.rootPath)

	var models = make(map[string]*Model)
	_ = filepath.Walk(r.rootPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if filepath.Ext(path) == ".json" {
				file, err := os.Open(path)
				if err != nil {
					log.Errorf(ctx, "Could not open: %s. %e", path, err)
					return nil
				}
				defer file.Close()

				b, err := ioutil.ReadAll(file)
				if err != nil {
					log.Errorf(ctx, "Could not read content from: %s. %e", path, err)
					return nil
				}

				model, err := r.unmarshal(path, string(b))
				if err != nil {
					log.Errorf(ctx, "Could not unmarshal content from: %s. %e", path, err)
					return nil
				}

				path = path[len(r.rootPath) : len(path)-5]
				path = strings.ToLower(path)
				path = strings.Replace(path, "\\", "/", -1)
				models[path] = model
				log.Infof(ctx, "Loaded %s", path)
			}
		}
		return nil
	})

	r.mux.Lock()
	defer r.mux.Unlock()
	r.Data = models
	return nil
}

func (r *RepositoryImpl) unmarshal(path string, str string) (*Model, error) {
	var raw pageData
	err := json.Unmarshal([]byte(str), &raw)
	if err != nil {
		return nil, errors.New("failed to unmarshal json: " + path + ". Reason: " + err.Error())
	}

	var content interface{}
	if raw.Type != "" {
		content, err = r.Types[raw.Type](raw.Content)
		if err != nil {
			return nil, errors.New("failed to unmarshal inner content: " + path + ". Reason: " + err.Error())
		}
	}

	return &Model{
		ID:        raw.ID,
		CreatedAt: raw.CreatedAt,
		View:      raw.View,
		Type:      raw.Type,
		Content:   content,
	}, nil
}

func (r *RepositoryImpl) OnEvent(ctx context.Context, e interface{}) error {
	if _, ok := e.(*event.Checkout); ok {
		if err := r.Reload(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (r *RepositoryImpl) Search(contentType string) []*SearchResult {
	r.mux.Lock()
	defer r.mux.Unlock()
	var result []*SearchResult
	for path, value := range r.Data {
		if value.Type == contentType {
			result = append(result, &SearchResult{path, value})
		}
	}
	return result
}

func (r *RepositoryImpl) Lookup(uuid string) *SearchResult {
	r.mux.Lock()
	defer r.mux.Unlock()
	for path, value := range r.Data {
		if value.ID == uuid {
			return &SearchResult{path, value}
		}
	}
	return nil
}

func (r *RepositoryImpl) GetAll() []*SearchResult {
	r.mux.Lock()
	defer r.mux.Unlock()
	var result []*SearchResult
	for path, value := range r.Data {
		result = append(result, &SearchResult{path, value})
	}
	return result
}

func NewRepository(bus *event.Bus, rootPath string) Repository {
	result := &RepositoryImpl{
		rootPath: rootPath,
		Types:    make(map[string]UnmarshalContentFunc),
	}
	bus.AddListener(result)
	return result
}
