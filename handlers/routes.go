package handlers

import "net/http"

func (h *handler) registerRoutes() {
	middlewareAuthentication := h.middlewareAuthentication

	h.registerRoute(http.MethodOptions, "^.+$", h.handleCORS())
	h.registerRoute(http.MethodPost, "^/users/login$", h.handleAuthentication())
	h.registerRoute(http.MethodPost, "^/users$", h.handleRegistration())
	h.registerRoute(http.MethodGet, "^/user$", middlewareAuthentication(h.handleGetCurrentUser(), true))
	h.registerRoute(http.MethodPut, "^/user$", middlewareAuthentication(h.handleUpdateUser(), true))
	h.registerRoute(http.MethodGet, "^/profiles/(?P<username>[\\w]+)$", middlewareAuthentication(h.handleGetProfile(), false))
	h.registerRoute(http.MethodPost, "^/profiles/(?P<username>[\\w]+)/follow$", middlewareAuthentication(h.handleFollowUser(), true))
	h.registerRoute(http.MethodDelete, "^/profiles/(?P<username>[\\w]+)/follow$", middlewareAuthentication(h.handleUnfollowUser(), true))
	h.registerRoute(http.MethodGet, "^/articles$", middlewareAuthentication(h.handleListArticles(), false))
	h.registerRoute(http.MethodGet, "^/articles/feed$", middlewareAuthentication(h.handleFeedArticles(), true))
	h.registerRoute(http.MethodGet, "^/articles/(?P<slug>[\\w-]+)$", h.handleGetArticle())
	h.registerRoute(http.MethodPost, "^/articles$", middlewareAuthentication(h.handleCreateArticle(), true))
	h.registerRoute(http.MethodPut, "^/articles/(?P<slug>[\\w-]+)$", middlewareAuthentication(h.handleUpdateArticle(), true))
	h.registerRoute(http.MethodDelete, "^/articles/(?P<slug>[\\w-]+)$", middlewareAuthentication(h.handleDeleteArticle(), true))
	h.registerRoute(http.MethodPost, "^/articles/(?P<slug>[\\w-]+)/comments$", middlewareAuthentication(h.handleAddCommentsToAnArticle(), true))
	h.registerRoute(http.MethodGet, "^/articles/(?P<slug>[\\w-]+)/comments$", middlewareAuthentication(h.handleGetCommentsFromAnArticle(), false))
	h.registerRoute(http.MethodDelete, "^/articles/(?P<slug>[\\w-]+)/comments/(?P<id>[\\d]+)$", middlewareAuthentication(h.handleDeleteComment(), true))
	h.registerRoute(http.MethodPost, "^/articles/(?P<slug>[\\w-]+)/favorite$", middlewareAuthentication(h.handleFavoriteArticle(), true))
	h.registerRoute(http.MethodDelete, "^/articles/(?P<slug>[\\w-]+)/favorite$", middlewareAuthentication(h.handleUnfavoriteArticle(), true))
	h.registerRoute(http.MethodGet, "^/tags$", h.handleGetTags())
}
