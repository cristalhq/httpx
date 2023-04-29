package httpx

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

// Server for HTTP protocol.
type Server struct {
	srv *http.Server
	cfg *ServerConfig
}

// ServerConfig configures Server.
type ServerConfig struct {
	Addr    string
	Handler http.Handler

	NoHTTP2   bool
	TLSConfig *tls.Config

	// TODO: enable option
	HeartbeatPath    string
	HeartbeatHandler http.HandlerFunc

	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	MaxHeaderBytes    int
}

// Validate the config.
func (c *ServerConfig) Validate() error {
	if c.ReadTimeout == 0 {
		c.ReadTimeout = 30 * time.Second
	}
	if c.ReadHeaderTimeout == 0 {
		c.ReadHeaderTimeout = 5 * time.Second
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = 15 * time.Second
	}
	if c.IdleTimeout == 0 {
		c.IdleTimeout = 30 * time.Second
	}
	if c.MaxHeaderBytes == 0 {
		c.MaxHeaderBytes = 8 * 1024
	}
	return nil
}

// NewServer returns a new Server.
func NewServer(config *ServerConfig) (*Server, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	s := &Server{
		srv: &http.Server{
			Addr:              config.Addr,
			Handler:           config.Handler,
			ReadTimeout:       config.ReadTimeout,
			ReadHeaderTimeout: config.ReadHeaderTimeout,
			WriteTimeout:      config.WriteTimeout,
			IdleTimeout:       config.IdleTimeout,
			MaxHeaderBytes:    config.MaxHeaderBytes,
		},
	}

	if config.NoHTTP2 {
		s.srv.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler))
	}
	return s, nil
}

// Start the server with a given handler.
// Same as Run but allows to set handler later.
func (s *Server) Start(ctx context.Context, h http.Handler) error {
	s.srv.Handler = h
	return s.Run(ctx)
}

// Run starts the server.
func (s *Server) Run(ctx context.Context) error {
	if s.srv.Handler == nil {
		panic("handler is nil")
	}

	s.srv.BaseContext = func(net.Listener) context.Context {
		return ctx
	}

	errCh := make(chan error)
	go func() {
		errCh <- s.srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		return s.shutdown()

	case err := <-errCh:
		return err
	}
}

func (s *Server) shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return s.srv.Shutdown(ctx)
}

func defaultHeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("."))
}
