package webserver

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	"net/http"
)

type RESTEndpoint struct {
	urlpath string
	verb    string
}

type WebServer struct {
	Router        chi.Router
	Handlers      map[RESTEndpoint]http.HandlerFunc
	WebServerPort string
	tracer        *zipkin.Tracer
}

func NewWebServer(serverPort string, tracer *zipkin.Tracer) *WebServer {
	if string(serverPort[0]) != ":" {
		serverPort = ":" + serverPort
	}
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[RESTEndpoint]http.HandlerFunc),
		WebServerPort: serverPort,
		tracer:        tracer,
	}
}

func (s *WebServer) AddHandler(urlpath string, verb string, handler http.HandlerFunc) {
	s.Handlers[RESTEndpoint{
		urlpath: urlpath,
		verb:    verb,
	}] = handler
}

// loop through the handlers and add them to the router
// register middeleware logger
// start the server
func (s *WebServer) Start() error {
	s.Router.Use(middleware.Logger)
	// create global zipkin http server middleware
	serverMiddleware := zipkinhttp.NewServerMiddleware(
		s.tracer, zipkinhttp.TagResponseSize(true),
	)
	s.Router.Use(serverMiddleware)

	for restEndpointInfo, handler := range s.Handlers {
		urlpath := restEndpointInfo.urlpath
		switch verb := restEndpointInfo.verb; verb {
		case http.MethodGet:
			s.Router.Get(urlpath, handler)
		case http.MethodPost:
			s.Router.Post(urlpath, handler)
		case http.MethodPut:
			s.Router.Put(urlpath, handler)
		case http.MethodPatch:
			s.Router.Patch(urlpath, handler)
		case http.MethodDelete:
			s.Router.Delete(urlpath, handler)
		default:
			return errors.New("invalid HTTP Verb")
		}
	}

	http.ListenAndServe(s.WebServerPort, s.Router)
	return nil
}
