package handlers

import (
	"context"
	"github.com/nasermirzaei89/realworld-go/models"
	"net/http"
	"regexp"
)

type Handler interface {
	http.Handler
}

type handler struct {
	userRepo    models.UserRepository
	articleRepo models.ArticleRepository
	routes      []route
	secret      []byte
}

type route struct {
	Method      string
	Pattern     regexp.Regexp
	HandlerFunc http.HandlerFunc
}

func NewHandler(userRepo models.UserRepository, articleRepo models.ArticleRepository, secret []byte) Handler {
	h := handler{
		userRepo:    userRepo,
		articleRepo: articleRepo,
		secret:      secret,
	}

	h.registerRoutes()

	return &h
}

func (h *handler) registerRoute(method, pattern string, handler http.HandlerFunc) {
	h.routes = append(h.routes, route{
		Method:      method,
		Pattern:     *regexp.MustCompile(pattern),
		HandlerFunc: handler,
	})
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		if r.Method == route.Method && route.Pattern.MatchString(r.URL.Path) {
			names := route.Pattern.SubexpNames()
			values := route.Pattern.FindAllStringSubmatch(r.URL.Path, -1)
			if len(values) > 0 {
				for i, v := range values[0] {
					if names[i] != "" {
						r = r.WithContext(context.WithValue(r.Context(), names[i], v))
					}
				}
			}
			route.HandlerFunc(w, r)
			return
		}
	}

	http.NotFoundHandler().ServeHTTP(w, r)
}
