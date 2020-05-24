package cms

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/westcoastcode-se/gocms/pkg/cache"
	"github.com/westcoastcode-se/gocms/pkg/config"
	"github.com/westcoastcode-se/gocms/pkg/content"
	"github.com/westcoastcode-se/gocms/pkg/event"
	. "github.com/westcoastcode-se/gocms/pkg/middleware"
	"github.com/westcoastcode-se/gocms/pkg/render"
	"github.com/westcoastcode-se/gocms/pkg/render/cached"
	"github.com/westcoastcode-se/gocms/pkg/render/immediate"
	"github.com/westcoastcode-se/gocms/pkg/security"
	"log"
	"net/http"
	"strings"
	"time"
)

type RequestContext struct {
	User     *security.User
	Response http.ResponseWriter
	Request  *http.Request
}

type RequestHandler interface {
	ServeHTTP(ctx *RequestContext)
}

type FileHandler struct {
	Prefix  string
	Handler http.Handler
}

type Server struct {
	// Bus for sending events to listeners
	Bus *event.Bus

	// Service used for managing user login.
	// Override this value if you want a custom service for managing a users security:
	//  server.SecurityService = custom.NewCustomerLoginService()
	SecurityService security.LoginService

	// LoginService used for generating a token based on the supplied user (and back to the user)
	// Override this value if you want to alter how the server tokenize a user before sending it to the client:
	//  server.Tokenizer = custom.NewTokenizer()
	Tokenizer security.Tokenizer

	// Controller for where content is located
	ContentController content.Controller

	// Repository where content can be found
	ContentRepository content.Repository

	// Handler for static files
	FileHandler FileHandler

	// Cache used for rendered pages
	PageCache cache.Pages

	// Used for figuring what parts of the web requires what user roles
	ACL security.ACL

	// Database containing all templates used when rendering the actual pages
	TemplateDatabase render.TemplateDatabase

	// Factory used when creating contexts used by the rendering framework
	// Override this if you want to add custom functions
	ContextFactory render.ContextFactory

	// Handlers
	Handlers map[string]RequestHandler

	config config.Config
	server http.Server
}

func (s *Server) ServeTemplate(rw http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path

	ctx := &RequestContext{User: r.Context().Value(security.SessionKey).(*security.User), Response: rw, Request: r}

	if s.handleBuiltIn(ctx) {
		return
	}
	if s.handleExtended(ctx) {
		return
	}

	Cache(s.PageCache, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		model, pageNotFound := s.ContentRepository.FindByPath(uri)

		// Create a new context used when rendering the template
		renderCtx := s.ContextFactory.Create(r, model)

		// Fetch a template that matches the view and url
		t, err := renderCtx.GetTemplate()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			log.Print(err)
			return
		}

		if pageNotFound != nil {
			rw.WriteHeader(http.StatusNotFound)
		} else {
			rw.WriteHeader(http.StatusOK)
		}

		err = t.ExecuteTemplate(rw, "index.html", model)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			log.Print(err)
		}
	})).ServeHTTP(rw, r)
}

// Figure out the IP address from the incoming request
func getIpAddress(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	start := time.Now()

	if r.URL.Path == "/" {
		r.URL.Path = "/index"
	}

	requestId := uuid.New().String()
	r = r.WithContext(context.WithValue(r.Context(), ContextKeyRequestID, requestId))
	token, err := r.Cookie(security.SessionKey)
	var user *security.User
	if err == nil && token.Value != "" {
		user, err = s.Tokenizer.TokenToUser(token.Value)
		if err != nil {
			log.Printf("User token could not be loaded. Reason %e\n", err)
		}
	}
	if user == nil {
		user = security.NotLoggedInUser
	}
	r = r.WithContext(context.WithValue(r.Context(), security.SessionKey, user))

	defer func() {
		diff := time.Since(start)
		log.Printf("")
		LogFromRequest(r).WithFields(logrus.Fields{
			"uri":     r.RequestURI,
			"method":  r.Method,
			"remote":  getIpAddress(r),
			"elapsed": diff,
		}).Info()
	}()

	Authorize(s.ACL,
		http.HandlerFunc(s.ServeTemplate)).ServeHTTP(rw, r)
}

func (s *Server) handleBuiltIn(ctx *RequestContext) bool {
	r := ctx.Request
	rw := ctx.Response
	uri := r.URL.Path
	if strings.HasPrefix(uri, "/api/v1") {
		uri = uri[len("/api/v1"):]
		if strings.HasPrefix(uri, "/authenticate") {
			uri = uri[len("/authenticate"):]
			if uri == "/login" {
				if r.Method != http.MethodPost {
					returnMethodNotAllowed(rw)
					return true
				}
				login(s.SecurityService, s.Tokenizer, ctx)
				return true
			} else if uri == "/logout" {
				if r.Method != http.MethodGet {
					returnMethodNotAllowed(rw)
					return true
				}
				logout(ctx)
				return true
			}
		} else if strings.HasPrefix(uri, "/checkout") {
			if r.Method != http.MethodPost {
				returnMethodNotAllowed(rw)
				return true
			}
			checkout(s.ContentController, ctx)
			return true
		} else if strings.HasPrefix(uri, "/pages") {
			if r.Method != http.MethodGet {
				returnMethodNotAllowed(rw)
				return true
			}
			getPages(s.ContentRepository, ctx)
			return true
		}
	}

	if strings.HasPrefix(uri, s.FileHandler.Prefix) {
		s.FileHandler.Handler.ServeHTTP(rw, r)
		return true
	}

	return false
}

func (s *Server) handleExtended(ctx *RequestContext) bool {
	uri := ctx.Request.URL.Path
	for key, value := range s.Handlers {
		if strings.HasPrefix(uri, key) {
			value.ServeHTTP(ctx)
			return true
		}
	}
	return false
}

func (s *Server) Handle(prefix string, handler RequestHandler) {
	s.Handlers[prefix] = handler
}

func (s *Server) ListenAndServe() error {
	err := s.ContentRepository.Reload()
	if err != nil {
		panic(err)
	}

	log.Printf("Listening for connections on %s\n", s.server.Addr)
	return s.server.ListenAndServe()
}

type wrapper struct {
	fn http.Handler
}

func (w *wrapper) ServeHTTP(ctx *RequestContext) {
	w.fn.ServeHTTP(ctx.Response, ctx.Request)
}

func Handler(fn http.Handler) RequestHandler {
	return &wrapper{fn: fn}
}

// Create a new CMS server by using the supplied configuration
func NewServer(config *config.Config) *Server {
	bus := event.NewBus()

	var pageCache cache.Pages
	if config.Author {
		pageCache = cache.NewNoCaching()
	} else {
		pageCache = cache.NewPermanentCache(bus, config.CacheDatabasePath)
	}

	contentRepository := content.NewRepository(bus, config.ContentDirectory+"/pages")
	var templateDatabase render.TemplateDatabase
	if config.Author {
		templateDatabase = immediate.NewFileSystemTemplateDatabase(config.ContentDirectory + "/templates")
	} else {
		templateDatabase = cached.NewDatabase(bus, config.ContentDirectory+"/templates")
	}
	contextFactory := &render.DefaultContextFactory{
		ContentRepository: contentRepository,
		TemplateDatabase:  templateDatabase,
		Config:            *config,
	}

	result := &Server{
		Bus:               bus,
		SecurityService:   security.NewLoginService(bus, config.UserDatabasePath),
		Tokenizer:         security.NewAsymmetricTokenizer(config.PublicKeyPath, config.PrivateKeyPath),
		ContentController: content.NewGitController(bus, config.ContentDirectory),
		ContentRepository: contentRepository,
		FileHandler: FileHandler{
			Prefix:  "/assets",
			Handler: http.FileServer(NewSecureFileSystem(config.ContentDirectory)),
		},
		PageCache:        pageCache,
		ACL:              security.NewFileBasedACL(bus, config.ACLDatabasePath),
		TemplateDatabase: templateDatabase,
		ContextFactory:   contextFactory,
		config:           *config,
		server: http.Server{
			Addr:         config.Server.ListenAddr,
			ReadTimeout:  time.Second * config.Server.ReadTimeout,
			WriteTimeout: time.Second * config.Server.WriteTimeout,
			IdleTimeout:  time.Second * config.Server.IdleTimeout,
		},
	}
	result.server.Handler = result
	return result
}
