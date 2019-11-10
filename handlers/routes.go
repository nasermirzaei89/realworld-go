package handlers

import "net/http"

func (h *handler) registerRoutes() {
	middlewareAuthentication := h.middlewareAuthentication

	h.registerRoute(http.MethodOptions, "^.+$", h.handleCORS())
	h.registerRoute(http.MethodPost, "^/api/users/login$", h.handleAuthentication())
	h.registerRoute(http.MethodPost, "^/api/users$", h.handleRegistration())
	h.registerRoute(http.MethodGet, "^/api/user$", middlewareAuthentication(h.handleGetCurrentUser(), true))
	h.registerRoute(http.MethodPut, "^/api/user$", middlewareAuthentication(h.handleUpdateUser(), true))
	h.registerRoute(http.MethodGet, "^/api/profiles/(?P<username>[\\w]+)$", middlewareAuthentication(h.handleGetProfile(), false))
	h.registerRoute(http.MethodPost, "^/api/profiles/(?P<username>[\\w]+)/follow$", middlewareAuthentication(h.handleFollowUser(), true))
	h.registerRoute(http.MethodDelete, "^/api/profiles/(?P<username>[\\w]+)/follow$", middlewareAuthentication(h.handleUnfollowUser(), true))
	h.registerRoute(http.MethodGet, "^/api/articles$", middlewareAuthentication(h.handleListArticles(), false))
	h.registerRoute(http.MethodGet, "^/api/articles/feed$", middlewareAuthentication(h.handleFeedArticles(), true))
	h.registerRoute(http.MethodGet, "^/api/articles/(?P<slug>[\\w]+)$", h.handleGetArticle())
	h.registerRoute(http.MethodPost, "^/api/articles$", middlewareAuthentication(h.handleCreateArticle(), true))
	h.registerRoute(http.MethodPut, "^/api/articles/(?P<slug>[\\w]+)$", middlewareAuthentication(h.handleUpdateArticle(), true))
	h.registerRoute(http.MethodDelete, "^/api/articles/(?P<slug>[\\w]+)$", middlewareAuthentication(h.handleDeleteArticle(), true))
	h.registerRoute(http.MethodPost, "^/api/articles/(?P<slug>[\\w]+)/comments$", middlewareAuthentication(h.handleAddCommentsToAnArticle(), true))
	h.registerRoute(http.MethodGet, "^/api/articles/(?P<slug>[\\w]+)/comments$", middlewareAuthentication(h.handleGetCommentsFromAnArticle(), false))
	h.registerRoute(http.MethodDelete, "^/api/articles/(?P<slug>[\\w]+)/comments/(?P<id>[\\w]+)$", middlewareAuthentication(h.handleDeleteComment(), true))
	h.registerRoute(http.MethodPost, "^/api/articles/(?P<slug>[\\w]+)/favorite$", middlewareAuthentication(h.handleFavoriteArticle(), true))
	h.registerRoute(http.MethodDelete, "^/api/articles/(?P<slug>[\\w]+)/favorite$", middlewareAuthentication(h.handleUnfavoriteArticle(), true))
	h.registerRoute(http.MethodGet, "^/api/tags$", h.handleGetTags())
}
