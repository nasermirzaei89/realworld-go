package handlers

import "net/http"

func (h *handler) registerRoutes() {
	middlewareAuthentication := h.middlewareAuthentication

	h.registerRoute(http.MethodPost, "^/api/users/login$", h.handleAuthentication())
	h.registerRoute(http.MethodPost, "^/api/users$", h.handleRegistration())
	h.registerRoute(http.MethodGet, "^/api/user$", middlewareAuthentication(h.handleGetCurrentUser()))
	h.registerRoute(http.MethodPut, "^/api/user$", middlewareAuthentication(h.handleUpdateUser()))
	h.registerRoute(http.MethodGet, "^/api/profiles/(?P<username>[\\w]+)$", h.handleGetProfile())
	h.registerRoute(http.MethodPost, "^/api/profiles/(?P<username>[\\w]+)/follow$", middlewareAuthentication(h.handleFollowUser()))
	h.registerRoute(http.MethodDelete, "^/api/profiles/(?P<username>[\\w]+)/follow$", middlewareAuthentication(h.handleUnfollowUser()))
	h.registerRoute(http.MethodGet, "^/api/articles$", h.handleListArticles())
	h.registerRoute(http.MethodGet, "^/api/articles/feed$", middlewareAuthentication(h.handleFeedArticles()))
	h.registerRoute(http.MethodGet, "^/api/articles/(?P<slug>[\\w]+)$", h.handleGetArticle())
	h.registerRoute(http.MethodPost, "^/api/articles$", middlewareAuthentication(h.handleCreateArticle()))
	h.registerRoute(http.MethodPut, "^/api/articles/(?P<slug>[\\w]+)$", middlewareAuthentication(h.handleUpdateArticle()))
	h.registerRoute(http.MethodDelete, "^/api/articles/(?P<slug>[\\w]+)$", middlewareAuthentication(h.handleDeleteArticle()))
	h.registerRoute(http.MethodPost, "^/api/articles/(?P<slug>[\\w]+)/comments$", middlewareAuthentication(h.handleAddCommentsToAnArticle()))
	h.registerRoute(http.MethodGet, "^/api/articles/(?P<slug>[\\w]+)/comments$", h.handleGetCommentsFromAnArticle())
	h.registerRoute(http.MethodDelete, "^/api/articles/(?P<slug>[\\w]+)/comments/(?P<id>[\\w]+)$", middlewareAuthentication(h.handleDeleteComment()))
	h.registerRoute(http.MethodPost, "^/api/articles/(?P<slug>[\\w]+)/favorite$", middlewareAuthentication(h.handleFavoriteArticle()))
	h.registerRoute(http.MethodDelete, "^/api/articles/(?P<slug>[\\w]+)/favorite$", middlewareAuthentication(h.handleUnfavoriteArticle()))
	h.registerRoute(http.MethodGet, "^/api/tags$", h.handleGetTags())
}
