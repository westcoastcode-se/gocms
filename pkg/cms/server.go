package cms

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/westcoastcode-se/gocms/pkg/cache"
	"github.com/westcoastcode-se/gocms/pkg/config"
	"github.com/westcoastcode-se/gocms/pkg/content"
	"github.com/westcoastcode-se/gocms/pkg/event"
	"github.com/westcoastcode-se/gocms/pkg/log"
	. "github.com/westcoastcode-se/gocms/pkg/middleware"
	"github.com/westcoastcode-se/gocms/pkg/render"
	"github.com/westcoastcode-se/gocms/pkg/render/html"
	"github.com/westcoastcode-se/gocms/pkg/render/html/cached"
	"github.com/westcoastcode-se/gocms/pkg/render/html/immediate"
	"github.com/westcoastcode-se/gocms/pkg/security"
	"github.com/westcoastcode-se/gocms/pkg/security/acl"
	"github.com/westcoastcode-se/gocms/pkg/security/auth"
	"github.com/westcoastcode-se/gocms/pkg/security/jwt"
	"net/http"
	"net/url"
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
	SecurityService auth.LoginService

	// LoginService used for generating a token based on the supplied user (and back to the user)
	// Override this value if you want to alter how the server tokenize a user before sending it to the client:
	//  server.Tokenizer = custom.NewTokenizer()
	Tokenizer jwt.Tokenizer

	// Controller for where content is located
	ContentController content.Controller

	// Repository where content can be found
	ContentRepository content.Repository

	// Handler for static files
	FileHandler FileHandler

	// Cache used for rendered pages
	PageCache cache.Pages

	// Used for figuring what parts of the web requires what user roles
	ACL acl.Service

	// Container for template renderers. You can add custom renderers if you want by:
	//  server.TemplateRenderers.AddFactory(NewCustomTemplateFactory())
	//
	// There's a built-in renderer for html files based on the "html/template" package. It's possible to override
	// this if you add a custom template renderer factory with html as suffix.
	TemplateRenderers *render.TemplateRenderers

	// Handlers
	Handlers map[string]RequestHandler

	config config.Config
	server http.Server
}

func (s *Server) ServeTemplate(rw http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path

	ctx := &RequestContext{User: r.Context().Value(jwt.SessionKey).(*security.User), Response: rw, Request: r}

	if s.handleBuiltIn(ctx) {
		return
	}
	if s.handleExtended(ctx) {
		return
	}

	Cache(s.PageCache, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		model, pageNotFound := s.ContentRepository.FindByPath(uri)

		// Fetch a factory for the template renderer. TODO: Custom view
		renderFactory, err := s.TemplateRenderers.FindFactory("index.html")
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			log.LogFromRequest(r).Warn(err.Error())
			return
		}

		// Set http status
		if pageNotFound != nil {
			rw.WriteHeader(http.StatusNotFound)
		}

		renderer := renderFactory.NewRenderer(r)
		err = renderer.RenderView(rw, "index.html", model)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			log.LogFromRequest(r).Warn(err.Error())
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
	r = r.WithContext(log.SetRequestID(r.Context(), requestId))

	token, err := r.Cookie(jwt.SessionKey)
	var user *security.User
	if err == nil && token.Value != "" {
		user, err = s.Tokenizer.TokenToUser(token.Value)
		if err != nil {
			log.LogFromRequest(r).Warnf("User token could not be loaded. Reason %e", err)
		}
	}
	if user == nil {
		user = security.NotLoggedInUser
	}

	r = r.WithContext(context.WithValue(r.Context(), jwt.SessionKey, user))
	r = r.WithContext(log.SetUserName(r.Context(), user.Name))

	defer func() {
		diff := time.Since(start)
		log.LogFromRequest(r).WithFields(logrus.Fields{
			"uri":     r.RequestURI,
			"method":  r.Method,
			"remote":  getIpAddress(r),
			"elapsed": diff,
		}).Info()
	}()

	uri := r.URL.Path
	roles := s.ACL.GetRoles(uri)
	if !user.HasRoles(roles) {
		http.Redirect(rw, r, "/login?redirect="+url.QueryEscape(uri), http.StatusTemporaryRedirect)
		return
	}

	s.ServeTemplate(rw, r)
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
	err := s.ContentRepository.Reload(context.Background())
	if err != nil {
		panic(err)
	}

	log.Infof(context.Background(), "Listening for connections on %s", s.server.Addr)
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

	var templateDatabase html.TemplateDatabase
	if config.Author {
		templateDatabase = immediate.NewFileSystemTemplateDatabase(config.ContentDirectory + "/templates")
	} else {
		templateDatabase = cached.NewDatabase(bus, config.ContentDirectory+"/templates")
	}

	templateRenderers := render.NewTemplateRenderers()
	templateRenderers.AddFactory(".html", &html.TemplateRendererFactory{
		ContentRepository: contentRepository,
		TemplateDatabase:  templateDatabase,
		Config:            *config,
	})

	result := &Server{
		Bus:               bus,
		SecurityService:   auth.NewLoginService(bus, config.UserDatabasePath),
		Tokenizer:         jwt.NewAsymmetricTokenizer(config.PublicKeyPath, config.PrivateKeyPath),
		ContentController: content.NewGitController(bus, config.ContentDirectory),
		ContentRepository: contentRepository,
		FileHandler: FileHandler{
			Prefix:  "/assets",
			Handler: http.FileServer(NewSecureFileSystem(config.ContentDirectory)),
		},
		PageCache:         pageCache,
		ACL:               acl.NewFileBasedACL(bus, config.ACLDatabasePath),
		TemplateRenderers: templateRenderers,
		config:            *config,
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
