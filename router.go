package httpx

import (
	"net/http"
	"slices"
)

// Router for net/http.
type Router struct {
	mux         *http.ServeMux
	globalMw    []func(http.Handler) http.Handler
	routeMw     []func(http.Handler) http.Handler
	isSubRouter bool
}

// NewRouter creates a new [Router]
func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

// Use given middlewares for the current router or subrouter.
func (r *Router) Use(mws ...func(http.Handler) http.Handler) {
	if r.isSubRouter {
		r.routeMw = append(r.routeMw, mws...)
	} else {
		r.globalMw = append(r.globalMw, mws...)
	}
}

// Group creates a new subrouter inheriting global middlewares.
func (r *Router) Group(fn func(r *Router)) {
	subRouter := &Router{
		mux:         r.mux,
		routeMw:     slices.Clone(r.routeMw),
		isSubRouter: true,
	}
	fn(subRouter)
}

// HandleFunc add pattern to the router or subrouter.
func (r *Router) HandleFunc(pattern string, h http.HandlerFunc) {
	r.Handle(pattern, h)
}

// Handle add pattern to the router or subrouter.
func (r *Router) Handle(pattern string, h http.Handler) {
	for _, mw := range slices.Backward(r.routeMw) {
		h = mw(h)
	}
	r.mux.Handle(pattern, h)
}

// ServeHTTP implements [http.Handler]
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var h http.Handler = r.mux

	for _, mw := range slices.Backward(r.globalMw) {
		h = mw(h)
	}
	h.ServeHTTP(w, req)
}
