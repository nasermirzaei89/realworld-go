package handlers

import "net/http"

func (h *handler) middlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		panic("implement me")
	}
}
