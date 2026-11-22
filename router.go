package httpx

import (
	"net/http"
	"slices"
)

// Router wraps net/http routing with global and grouped middleware.
// Routers passed to Group are registration scopes only.
type Router struct {
	mux         *http.ServeMux
	handler     http.Handler
	globalMw    []func(http.Handler) http.Handler
	routeMw     []func(http.Handler) http.Handler
	isSubRouter bool
}

// NewRouter creates a new [Router].
func NewRouter() *Router {
	mux := http.NewServeMux()
	return &Router{
		mux:     mux,
		handler: mux,
	}
}

// ServeHTTP implements [http.Handler].
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.handler.ServeHTTP(w, req)
}

// Use adds middleware in execution order.
// On the root router, middleware applies to all requests.
// Within a group, middleware applies to subsequently registered routes.
func (r *Router) Use(mws ...func(http.Handler) http.Handler) {
	if len(mws) == 0 {
		return
	}

	if r.isSubRouter {
		r.routeMw = append(r.routeMw, mws...)
		return
	}

	r.globalMw = append(r.globalMw, mws...)

	var h http.Handler = r.mux
	for _, mw := range slices.Backward(r.globalMw) {
		h = mw(h)
	}
	r.handler = h
}

// Group calls fn with a registration scope sharing the root router's mux.
// The scope inherits a snapshot of its parent's route middleware.
// Middleware added within the scope does not affect its parent or siblings.
// Global middleware applies when requests are served through the root router.
func (r *Router) Group(fn func(r *Router)) {
	subRouter := &Router{
		mux:         r.mux,
		handler:     r.mux,
		routeMw:     slices.Clone(r.routeMw),
		isSubRouter: true,
	}
	fn(subRouter)
}

// HandleFunc registers h for pattern.
func (r *Router) HandleFunc(pattern string, h http.HandlerFunc) {
	if h == nil {
		panic("httpx: nil handler")
	}
	r.Handle(pattern, h)
}

// Handle registers h for pattern using the current route middleware.
func (r *Router) Handle(pattern string, h http.Handler) {
	if h == nil {
		panic("httpx: nil handler")
	}

	for _, mw := range slices.Backward(r.routeMw) {
		h = mw(h)
	}
	r.mux.Handle(pattern, h)
}
