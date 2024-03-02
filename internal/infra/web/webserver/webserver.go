package webserver

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/openzipkin/zipkin-go"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"net/http"
)

type RESTEndpoint struct {
	urlpath     string
	verb        string
	requestName string
}

type WebServer struct {
	Router        *mux.Router
	Handlers      map[RESTEndpoint]http.HandlerFunc
	WebServerPort string
	Tracer        *zipkin.Tracer
}

func NewWebServer(serverPort string) *WebServer {
	if string(serverPort[0]) != ":" {
		serverPort = ":" + serverPort
	}
	return &WebServer{
		Router:        mux.NewRouter(),
		Handlers:      make(map[RESTEndpoint]http.HandlerFunc),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(urlpath string, verb string, handler http.HandlerFunc, requestName string) {
	s.Handlers[RESTEndpoint{
		urlpath:     urlpath,
		verb:        verb,
		requestName: requestName,
	}] = handler
}

// loop through the handlers and add them to the router
// register middeleware logger
// start the server
func (s *WebServer) Start() error {
	s.Router.Use(otelmux.Middleware("zipkin-goexpert"))

	for restEndpointInfo, handler := range s.Handlers {
		urlpath := restEndpointInfo.urlpath

		switch verb := restEndpointInfo.verb; verb {
		case http.MethodGet:
			s.Router.HandleFunc(urlpath, handler).Methods(http.MethodGet)
		case http.MethodPost:
			s.Router.HandleFunc(urlpath, handler).Methods(http.MethodPost)
		case http.MethodPut:
			s.Router.HandleFunc(urlpath, handler).Methods(http.MethodPut)
		case http.MethodPatch:
			s.Router.HandleFunc(urlpath, handler).Methods(http.MethodPatch)
		case http.MethodDelete:
			s.Router.HandleFunc(urlpath, handler).Methods(http.MethodDelete)
		default:
			return errors.New("invalid HTTP Verb")
		}
	}
	log.Info().Msgf("Starting webserver at port %s", s.WebServerPort)
	err := http.ListenAndServe(s.WebServerPort, s.Router)
	return err
}
