package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nasermirzaei89/realworld-go/internal/models"
	uniqueID "github.com/nasermirzaei89/realworld-go/pkg/id"
	"github.com/nasermirzaei89/realworld-go/pkg/jwt"
	slugify "github.com/nasermirzaei89/realworld-go/pkg/slug"
	"net/http"
	"strconv"
	"time"
)

func (h *handler) handleCORS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *handler) handleAuthentication() http.HandlerFunc {
	type Request struct {
		User struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		} `json:"user"`
	}

	type Response UserResponse

	return func(w http.ResponseWriter, r *http.Request) {
		// get request body
		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnprocessableEntity)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"body": err.Error(),
				},
			})
			return
		}

		// find user by email
		user, err := h.userRepo.GetByEmail(req.User.Email)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "get user by email failed",
					"error":   err.Error(),
				},
			})
			return
		}

		// check password
		// TODO: should hash password
		if req.User.Password != user.Password {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "invalid password received",
				},
			})
			return
		}

		// success response
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Response{
			User: User{
				Email:    user.Email,
				Token:    user.Token,
				Username: user.Username,
				Bio:      user.Bio,
				Image:    user.Image,
			},
		})
	}
}

func (h *handler) handleRegistration() http.HandlerFunc {
	type Request struct {
		User struct {
			Email    string `json:"email"`
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"user"`
	}

	type Response UserResponse

	return func(w http.ResponseWriter, r *http.Request) {
		// get request body
		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnprocessableEntity)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"body": err.Error(),
				},
			})
			return
		}

		// find user by email
		_, err = h.userRepo.GetByEmail(req.User.Email)
		if err != nil && !errors.As(err, &models.UserByEmailNotFoundError{}) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "get user by email failed",
					"error":   err.Error(),
				},
			})
			return
		}

		if err == nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnprocessableEntity)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "email already taken",
				},
			})
			return
		}

		// find user by username
		_, err = h.userRepo.GetByUsername(req.User.Username)
		if err != nil && !errors.As(err, &models.UserByUsernameNotFoundError{}) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "get user by username failed",
					"error":   err.Error(),
				},
			})
			return
		}

		if err == nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnprocessableEntity)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "username already taken",
				},
			})
			return
		}

		// generate token with no expiration
		userID := h.userRepo.NewID()
		token := jwt.New()
		token.SetSubject(strconv.Itoa(userID))
		tokenStr, err := jwt.Sign(token, h.secret)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "error on sign jwt token",
					"error":   err.Error(),
				},
			})
			return
		}

		// create user
		user := models.User{
			ID:        userID,
			Email:     req.User.Email, // TODO: validate email
			Token:     tokenStr,
			Username:  req.User.Username,
			Password:  req.User.Password, // TODO: should hash password
			Bio:       "",
			Image:     "",
			Followers: map[int]bool{},
		}

		err = h.userRepo.Add(user)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "create user failed",
					"error":   err.Error(),
				},
			})
			return
		}

		// success response
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Response{
			User: User{
				Email:    user.Email,
				Token:    user.Token,
				Username: user.Username,
				Bio:      user.Bio,
				Image:    user.Image,
			},
		})
	}
}

func (h *handler) handleGetCurrentUser() http.HandlerFunc {
	type Response UserResponse

	return func(w http.ResponseWriter, r *http.Request) {
		// get current user
		currentUser := r.Context().Value(currentUserCtx).(*models.User)

		// success response
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Response{
			User: User{
				Email:    currentUser.Email,
				Token:    currentUser.Token,
				Username: currentUser.Username,
				Bio:      currentUser.Bio,
				Image:    currentUser.Image,
			},
		})
	}
}

func (h *handler) handleUpdateUser() http.HandlerFunc {
	type Request struct {
		User struct {
			Email    *string `json:"email"`
			Username *string `json:"username"`
			Password *string `json:"password"`
			Image    *string `json:"image"`
			Bio      *string `json:"bio"`
		} `json:"user"`
	}

	type Response UserResponse

	return func(w http.ResponseWriter, r *http.Request) {
		// get current user
		currentUser := r.Context().Value(currentUserCtx).(*models.User)

		// get request body
		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnprocessableEntity)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"body": err.Error(),
				},
			})
			return
		}

		if req.User.Email != nil {
			// check email
			exists, err := h.userRepo.GetByEmail(*req.User.Email)
			if err != nil && !errors.As(err, &models.UserByEmailNotFoundError{}) {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": "get user by email failed",
						"error":   err.Error(),
					},
				})
				return
			}

			if err == nil && exists.ID != currentUser.ID {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusUnprocessableEntity)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": "email already taken",
					},
				})
				return
			}

			// TODO: validate email
			currentUser.Email = *req.User.Email
		}

		if req.User.Username != nil {
			// check username
			exists, err := h.userRepo.GetByUsername(*req.User.Username)
			if err != nil && !errors.As(err, &models.UserByUsernameNotFoundError{}) {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": "get user by username failed",
						"error":   err.Error(),
					},
				})
				return
			}

			if err == nil && exists.ID != currentUser.ID {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusUnprocessableEntity)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": "username already taken",
					},
				})
				return
			}

			currentUser.Username = *req.User.Username
		}

		if req.User.Password != nil {
			currentUser.Password = *req.User.Password
		}

		if req.User.Image != nil {
			currentUser.Image = *req.User.Image
		}

		if req.User.Bio != nil {
			currentUser.Bio = *req.User.Bio
		}

		// update user
		err = h.userRepo.UpdateByID(currentUser.ID, *currentUser)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "update user failed",
					"error":   err.Error(),
				},
			})
			return
		}

		// success response
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Response{
			User: User{
				Email:    currentUser.Email,
				Token:    currentUser.Token,
				Username: currentUser.Username,
				Bio:      currentUser.Bio,
				Image:    currentUser.Image,
			},
		})
	}
}

func (h *handler) handleGetProfile() http.HandlerFunc {
	type Response ProfileResponse

	return func(w http.ResponseWriter, r *http.Request) {
		// get current user if exists
		var currentUser *models.User
		if iCurrentUser := r.Context().Value(currentUserCtx); iCurrentUser != nil {
			currentUser = iCurrentUser.(*models.User)
		}

		// get user by username
		user, err := h.userRepo.GetByUsername(r.Context().Value("username").(string))
		if err != nil {
			if errors.As(err, &models.UserByUsernameNotFoundError{}) {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": "user not found",
						"error":   err.Error(),
					},
				})
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "get user by username failed",
					"error":   err.Error(),
				},
			})
			return
		}

		following := false
		if currentUser != nil {
			if f, ok := user.Followers[currentUser.ID]; f && ok {
				following = true
			}
		}

		// success response
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Response{
			Profile: Profile{
				Username:  user.Username,
				Bio:       user.Bio,
				Image:     user.Image,
				Following: following,
			},
		})
	}
}

func (h *handler) handleFollowUser() http.HandlerFunc {
	type Response ProfileResponse

	return func(w http.ResponseWriter, r *http.Request) {
		// get current user
		currentUser := r.Context().Value(currentUserCtx).(*models.User)

		// get user by username
		user, err := h.userRepo.GetByUsername(r.Context().Value("username").(string))
		if err != nil {
			if errors.As(err, &models.UserByUsernameNotFoundError{}) {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": "user not found",
						"error":   err.Error(),
					},
				})
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "get user by username failed",
					"error":   err.Error(),
				},
			})
			return
		}

		user.Followers[currentUser.ID] = true

		// update user
		err = h.userRepo.UpdateByID(currentUser.ID, *currentUser)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "update user failed",
					"error":   err.Error(),
				},
			})
			return
		}

		// success response
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Response{
			Profile: Profile{
				Username:  user.Username,
				Bio:       user.Bio,
				Image:     user.Image,
				Following: true,
			},
		})
	}
}

func (h *handler) handleUnfollowUser() http.HandlerFunc {
	type Response ProfileResponse

	return func(w http.ResponseWriter, r *http.Request) {
		// get current user
		currentUser := r.Context().Value(currentUserCtx).(*models.User)

		// get user by username
		user, err := h.userRepo.GetByUsername(r.Context().Value("username").(string))
		if err != nil {
			if errors.As(err, &models.UserByUsernameNotFoundError{}) {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": "user not found",
						"error":   err.Error(),
					},
				})
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "get user by username failed",
					"error":   err.Error(),
				},
			})
			return
		}

		delete(user.Followers, currentUser.ID)

		// update user
		err = h.userRepo.UpdateByID(currentUser.ID, *currentUser)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "update user failed",
					"error":   err.Error(),
				},
			})
			return
		}

		// success response
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Response{
			Profile: Profile{
				Username:  user.Username,
				Bio:       user.Bio,
				Image:     user.Image,
				Following: false,
			},
		})
	}
}

func (h *handler) handleListArticles() http.HandlerFunc {
	type Response MultipleArticlesResponse

	return func(w http.ResponseWriter, r *http.Request) {
		// get current user if exists
		var currentUser *models.User
		if iCurrentUser := r.Context().Value(currentUserCtx); iCurrentUser != nil {
			currentUser = iCurrentUser.(*models.User)
		}

		query := r.URL.Query()

		var (
			filters []models.ArticleFilter
			offset  = 0
			limit   = 20
		)

		for k, vv := range query {
			for _, v := range vv {
				switch k {
				case "tag":
					filters = append(filters, models.FilterArticlesByTag(v))
				case "author":
					user, err := h.userRepo.GetByUsername(v)
					if err != nil {
						if errors.As(err, &models.UserByUsernameNotFoundError{}) {
							user = &models.User{Username: v}
						} else {
							w.Header().Set("Content-Type", "application/json; charset=utf-8")
							w.WriteHeader(http.StatusInternalServerError)
							_ = json.NewEncoder(w).Encode(ErrorResponse{
								Errors: map[string]interface{}{
									"message": "get user by username failed",
									"error":   err.Error(),
								},
							})
							return
						}
					}
					filters = append(filters, models.FilterArticlesByAuthor(*user))
				case "favorited":
					user, err := h.userRepo.GetByUsername(v)
					if err != nil {
						if errors.As(err, &models.UserByUsernameNotFoundError{}) {
							user = &models.User{Username: v}
						} else {
							w.Header().Set("Content-Type", "application/json; charset=utf-8")
							w.WriteHeader(http.StatusInternalServerError)
							_ = json.NewEncoder(w).Encode(ErrorResponse{
								Errors: map[string]interface{}{
									"message": "get user by username failed",
									"error":   err.Error(),
								},
							})
							return
						}
					}
					filters = append(filters, models.FilterArticlesByFavorite(*user))
				case "offset":
					var err error
					offset, err = strconv.Atoi(v)
					if err != nil {
						w.Header().Set("Content-Type", "application/json; charset=utf-8")
						w.WriteHeader(http.StatusUnprocessableEntity)
						_ = json.NewEncoder(w).Encode(ErrorResponse{
							Errors: map[string]interface{}{
								"message": "invalid offset received",
								"error":   err.Error(),
							},
						})
						return
					}
				case "limit":
					var err error
					limit, err = strconv.Atoi(v)
					if err != nil {
						w.Header().Set("Content-Type", "application/json; charset=utf-8")
						w.WriteHeader(http.StatusUnprocessableEntity)
						_ = json.NewEncoder(w).Encode(ErrorResponse{
							Errors: map[string]interface{}{
								"message": "invalid limit received",
								"error":   err.Error(),
							},
						})
						return
					}
				default:
					// drop filter
				}
			}
		}

		res, total, err := h.articleRepo.List(offset, limit, filters...)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "list article failed",
					"error":   err.Error(),
				},
			})
			return
		}

		articles := make([]Article, len(res))

		for i := range res {
			favorited := false
			if currentUser != nil {
				if f, ok := res[i].Favorites[currentUser.ID]; f && ok {
					favorited = true
				}
			}

			user, err := h.userRepo.GetByID(res[i].AuthorID)
			if err != nil {
				if errors.As(err, &models.UserByIDNotFoundError{}) {
					w.Header().Set("Content-Type", "application/json; charset=utf-8")
					w.WriteHeader(http.StatusNotFound)
					_ = json.NewEncoder(w).Encode(ErrorResponse{
						Errors: map[string]interface{}{
							"message": "user not found",
							"error":   err.Error(),
						},
					})
					return
				}

				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": "get user by id failed",
						"error":   err.Error(),
					},
				})
				return
			}

			following := false
			if currentUser != nil {
				if f, ok := user.Followers[currentUser.ID]; f && ok {
					following = true
				}
			}

			articles[i] = Article{
				Slug:           res[i].Slug,
				Title:          res[i].Title,
				Description:    res[i].Description,
				Body:           res[i].Body,
				TagList:        res[i].Tags,
				CreatedAt:      res[i].CreatedAt.UTC().Format(dateLayout),
				UpdatedAt:      res[i].UpdatedAt.UTC().Format(dateLayout),
				Favorited:      favorited,
				FavoritesCount: len(res[i].Favorites),
				Author: Author{
					Username:  user.Username,
					Bio:       user.Bio,
					Image:     user.Image,
					Following: following,
				},
			}
		}

		// success response
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Response{
			Articles:      articles,
			ArticlesCount: total,
		})
	}
}

func (h *handler) handleFeedArticles() http.HandlerFunc {
	type Response MultipleArticlesResponse

	return func(w http.ResponseWriter, r *http.Request) {
		// get current user
		currentUser := r.Context().Value(currentUserCtx).(*models.User)

		query := r.URL.Query()

		var (
			filters []models.ArticleFilter
			offset  = 0
			limit   = 20
		)

		for k, vv := range query {
			for _, v := range vv {
				switch k {
				case "offset":
					var err error
					offset, err = strconv.Atoi(v)
					if err != nil {
						w.Header().Set("Content-Type", "application/json; charset=utf-8")
						w.WriteHeader(http.StatusUnprocessableEntity)
						_ = json.NewEncoder(w).Encode(ErrorResponse{
							Errors: map[string]interface{}{
								"message": "invalid offset received",
								"error":   err.Error(),
							},
						})
						return
					}
				case "limit":
					var err error
					limit, err = strconv.Atoi(v)
					if err != nil {
						w.Header().Set("Content-Type", "application/json; charset=utf-8")
						w.WriteHeader(http.StatusUnprocessableEntity)
						_ = json.NewEncoder(w).Encode(ErrorResponse{
							Errors: map[string]interface{}{
								"message": "invalid limit received",
								"error":   err.Error(),
							},
						})
						return
					}
				default:
					// drop filter
				}
			}
		}

		// get followee
		users, err := h.userRepo.ListByFollowedBy(currentUser.ID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "error on get user followee",
					"error":   err.Error(),
				},
			})
			return
		}

		filters = append(filters, models.FilterArticlesByAuthors(users...))

		res, total, err := h.articleRepo.List(offset, limit, filters...)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "list article failed",
					"error":   err.Error(),
				},
			})
			return
		}

		articles := make([]Article, len(res))

		for i := range res {
			favorited := false
			if currentUser != nil {
				if f, ok := res[i].Favorites[currentUser.ID]; f && ok {
					favorited = true
				}
			}

			user, err := h.userRepo.GetByID(res[i].AuthorID)
			if err != nil {
				if errors.As(err, &models.UserByIDNotFoundError{}) {
					w.Header().Set("Content-Type", "application/json; charset=utf-8")
					w.WriteHeader(http.StatusNotFound)
					_ = json.NewEncoder(w).Encode(ErrorResponse{
						Errors: map[string]interface{}{
							"message": "user not found",
							"error":   err.Error(),
						},
					})
					return
				}

				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": "get user by id failed",
						"error":   err.Error(),
					},
				})
				return
			}

			following := false
			if currentUser != nil {
				if f, ok := user.Followers[currentUser.ID]; f && ok {
					following = true
				}
			}

			articles[i] = Article{
				Slug:           res[i].Slug,
				Title:          res[i].Title,
				Description:    res[i].Description,
				Body:           res[i].Body,
				TagList:        res[i].Tags,
				CreatedAt:      res[i].CreatedAt.UTC().Format(dateLayout),
				UpdatedAt:      res[i].UpdatedAt.UTC().Format(dateLayout),
				Favorited:      favorited,
				FavoritesCount: len(res[i].Favorites),
				Author: Author{
					Username:  user.Username,
					Bio:       user.Bio,
					Image:     user.Image,
					Following: following,
				},
			}
		}

		// success response
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Response{
			Articles:      articles,
			ArticlesCount: total,
		})
	}
}

func (h *handler) handleGetArticle() http.HandlerFunc {
	type Response SingleArticleResponse

	return func(w http.ResponseWriter, r *http.Request) {
		// get current user if exists
		var currentUser *models.User
		if iCurrentUser := r.Context().Value(currentUserCtx); iCurrentUser != nil {
			currentUser = iCurrentUser.(*models.User)
		}

		// get params
		slug := r.Context().Value("slug").(string)

		// find article by slug
		article, err := h.articleRepo.GetBySlug(slug)
		if err != nil {
			if errors.As(err, &models.ArticleBySlugNotFoundError{}) {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": fmt.Sprintf("article with slug '%s' not found", slug),
						"error":   err.Error(),
					},
				})
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "get article by slug failed",
					"error":   err.Error(),
				},
			})
			return
		}

		// find user by id
		user, err := h.userRepo.GetByID(article.AuthorID)
		if err != nil {
			if errors.As(err, &models.UserByIDNotFoundError{}) {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": "author of article not found",
						"error":   err.Error(),
					},
				})
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "get author of article failed",
					"error":   err.Error(),
				},
			})
			return
		}

		favorited := false
		if currentUser != nil {
			if f, ok := article.Favorites[currentUser.ID]; f && ok {
				favorited = true
			}
		}

		following := false
		if currentUser != nil {
			if f, ok := user.Followers[currentUser.ID]; f && ok {
				following = true
			}
		}

		// success response
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Response{
			Article: Article{
				Slug:           article.Slug,
				Title:          article.Title,
				Description:    article.Description,
				Body:           article.Body,
				TagList:        article.Tags,
				CreatedAt:      article.CreatedAt.UTC().Format(dateLayout),
				UpdatedAt:      article.UpdatedAt.UTC().Format(dateLayout),
				Favorited:      favorited,
				FavoritesCount: len(article.Favorites),
				Author: Author{
					Username:  user.Username,
					Bio:       user.Bio,
					Image:     user.Image,
					Following: following,
				},
			},
		})
	}
}

func (h *handler) handleCreateArticle() http.HandlerFunc {
	type Request struct {
		Article struct {
			Title       string   `json:"title"`
			Description string   `json:"description"`
			Body        string   `json:"body"`
			TagList     []string `json:"tagList"`
		} `json:"article"`
	}

	type Response SingleArticleResponse

	return func(w http.ResponseWriter, r *http.Request) {
		// get current user
		currentUser := r.Context().Value(currentUserCtx).(*models.User)

		// get request body
		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnprocessableEntity)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"body": err.Error(),
				},
			})
			return
		}

		// create article
		article := models.Article{
			Slug:        fmt.Sprintf("%s-%s", slugify.Make(req.Article.Title), uniqueID.New(6)),
			Title:       req.Article.Title,
			Description: req.Article.Description,
			Body:        req.Article.Body,
			Tags:        append([]string{}, req.Article.TagList...),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			AuthorID:    currentUser.ID,
			Favorites:   make(map[int]bool),
			Comments:    make([]models.Comment, 0),
		}

		err = h.articleRepo.Add(article)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "error on create article",
					"body":    err.Error(),
				},
			})
			return
		}

		// success response
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(Response{
			Article: Article{
				Slug:           article.Slug,
				Title:          article.Title,
				Description:    article.Description,
				Body:           article.Body,
				TagList:        article.Tags,
				CreatedAt:      article.CreatedAt.UTC().Format(dateLayout),
				UpdatedAt:      article.UpdatedAt.UTC().Format(dateLayout),
				Favorited:      false,
				FavoritesCount: len(article.Favorites),
				Author: Author{
					Username:  currentUser.Username,
					Bio:       currentUser.Bio,
					Image:     currentUser.Image,
					Following: false,
				},
			},
		})
	}
}

func (h *handler) handleUpdateArticle() http.HandlerFunc {
	type Request struct {
		Article struct {
			Title       *string `json:"title"`
			Description *string `json:"description"`
			Body        *string `json:"body"`
		} `json:"article"`
	}

	type Response SingleArticleResponse

	return func(w http.ResponseWriter, r *http.Request) {
		// get current user
		currentUser := r.Context().Value(currentUserCtx).(*models.User)

		// get params
		slug := r.Context().Value("slug").(string)

		// find article by slug
		article, err := h.articleRepo.GetBySlug(slug)
		if err != nil {
			if errors.As(err, &models.ArticleBySlugNotFoundError{}) {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": fmt.Sprintf("article with slug '%s' not found", slug),
						"error":   err.Error(),
					},
				})
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "get article by slug failed",
					"error":   err.Error(),
				},
			})
			return
		}

		// check author
		if article.AuthorID != currentUser.ID {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "you are not author of this article",
				},
			})
			return
		}

		// get request body
		var req Request
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnprocessableEntity)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"body": err.Error(),
				},
			})
			return
		}

		// update fields
		if req.Article.Title != nil {
			article.Title = *req.Article.Title
			article.Slug = fmt.Sprintf("%s-%s", slugify.Make(*req.Article.Title), uniqueID.New(6))
		}

		if req.Article.Description != nil {
			article.Description = *req.Article.Description
		}

		if req.Article.Body != nil {
			article.Body = *req.Article.Body
		}

		// update article
		err = h.articleRepo.UpdateBySlug(slug, *article)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "error on update article",
					"body":    err.Error(),
				},
			})
			return
		}

		favorited := false
		if currentUser != nil {
			if f, ok := article.Favorites[currentUser.ID]; f && ok {
				favorited = true
			}
		}

		following := false
		if f, ok := currentUser.Followers[currentUser.ID]; f && ok {
			following = true
		}

		// success response
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Response{
			Article: Article{
				Slug:           article.Slug,
				Title:          article.Title,
				Description:    article.Description,
				Body:           article.Body,
				TagList:        article.Tags,
				CreatedAt:      article.CreatedAt.UTC().Format(dateLayout),
				UpdatedAt:      article.UpdatedAt.UTC().Format(dateLayout),
				Favorited:      favorited,
				FavoritesCount: len(article.Favorites),
				Author: Author{
					Username:  currentUser.Username,
					Bio:       currentUser.Bio,
					Image:     currentUser.Image,
					Following: following,
				},
			},
		})
	}
}

func (h *handler) handleDeleteArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get current user
		currentUser := r.Context().Value(currentUserCtx).(*models.User)

		// get params
		slug := r.Context().Value("slug").(string)

		// find article by slug
		article, err := h.articleRepo.GetBySlug(slug)
		if err != nil {
			if errors.As(err, &models.ArticleBySlugNotFoundError{}) {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": fmt.Sprintf("article with slug '%s' not found", slug),
						"error":   err.Error(),
					},
				})
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "get article by slug failed",
					"error":   err.Error(),
				},
			})
			return
		}

		// check author
		if article.AuthorID != currentUser.ID {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "you are not author of this article",
				},
			})
			return
		}

		// delete article
		err = h.articleRepo.DeleteBySlug(slug)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "error on update article",
					"body":    err.Error(),
				},
			})
			return
		}

		// success response
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *handler) handleAddCommentsToAnArticle() http.HandlerFunc {
	type Request struct {
		Comment struct {
			Body string `json:"body"`
		} `json:"comment"`
	}

	type Response SingleCommentResponse

	return func(w http.ResponseWriter, r *http.Request) {
		// get current user
		currentUser := r.Context().Value(currentUserCtx).(*models.User)

		// get params
		slug := r.Context().Value("slug").(string)

		// find article by slug
		article, err := h.articleRepo.GetBySlug(slug)
		if err != nil {
			if errors.As(err, &models.ArticleBySlugNotFoundError{}) {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": fmt.Sprintf("article with slug '%s' not found", slug),
						"error":   err.Error(),
					},
				})
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "get article by slug failed",
					"error":   err.Error(),
				},
			})
			return
		}

		// get request body
		var req Request
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnprocessableEntity)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"body": err.Error(),
				},
			})
			return
		}

		// create comment
		comment := models.Comment{
			ID:        h.articleRepo.NewCommentID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Body:      req.Comment.Body,
			AuthorID:  currentUser.ID,
		}

		article.Comments = append(article.Comments, comment)

		// update article
		err = h.articleRepo.UpdateBySlug(slug, *article)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "error on update article",
					"body":    err.Error(),
				},
			})
			return
		}

		following := false
		if f, ok := currentUser.Followers[currentUser.ID]; f && ok {
			following = true
		}

		// success response
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(Response{
			Comment: Comment{
				ID:        comment.ID,
				CreatedAt: comment.CreatedAt.UTC().Format(dateLayout),
				UpdatedAt: comment.UpdatedAt.UTC().Format(dateLayout),
				Body:      comment.Body,
				Author: Author{
					Username:  currentUser.Username,
					Bio:       currentUser.Bio,
					Image:     currentUser.Image,
					Following: following,
				},
			},
		})
	}
}

func (h *handler) handleGetCommentsFromAnArticle() http.HandlerFunc {
	type Response MultipleCommentsResponse

	return func(w http.ResponseWriter, r *http.Request) {
		// get current user if exists
		var currentUser *models.User
		if iCurrentUser := r.Context().Value(currentUserCtx); iCurrentUser != nil {
			currentUser = iCurrentUser.(*models.User)
		}

		// get params
		slug := r.Context().Value("slug").(string)

		// find article by slug
		article, err := h.articleRepo.GetBySlug(slug)
		if err != nil {
			if errors.As(err, &models.ArticleBySlugNotFoundError{}) {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": fmt.Sprintf("article with slug '%s' not found", slug),
						"error":   err.Error(),
					},
				})
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "get article by slug failed",
					"error":   err.Error(),
				},
			})
			return
		}

		comments := make([]Comment, len(article.Comments))
		for i := range article.Comments {
			author, err := h.userRepo.GetByID(article.Comments[i].AuthorID)
			if err != nil {
				if errors.As(err, &models.UserByIDNotFoundError{}) {
					w.Header().Set("Content-Type", "application/json; charset=utf-8")
					w.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(w).Encode(ErrorResponse{
						Errors: map[string]interface{}{
							"message": "author of comment not found",
							"error":   err.Error(),
						},
					})
					return
				}

				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": "get author of comment failed",
						"error":   err.Error(),
					},
				})
				return
			}

			following := false
			if currentUser != nil {
				if f, ok := author.Followers[currentUser.ID]; f && ok {
					following = true
				}
			}

			comments[i] = Comment{
				ID:        article.Comments[i].ID,
				CreatedAt: article.Comments[i].CreatedAt.UTC().Format(dateLayout),
				UpdatedAt: article.Comments[i].UpdatedAt.UTC().Format(dateLayout),
				Body:      article.Comments[i].Body,
				Author: Author{
					Username:  author.Username,
					Bio:       author.Bio,
					Image:     author.Image,
					Following: following,
				},
			}
		}

		// success response
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Response{
			Comments: comments,
		})
	}
}

func (h *handler) handleDeleteComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get current user
		currentUser := r.Context().Value(currentUserCtx).(*models.User)

		// get params
		slug := r.Context().Value("slug").(string)
		idStr := r.Context().Value("id").(string)
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnprocessableEntity)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "invalid comment id received",
					"error":   err.Error(),
				},
			})
			return
		}

		// find article by slug
		article, err := h.articleRepo.GetBySlug(slug)
		if err != nil {
			if errors.As(err, &models.ArticleBySlugNotFoundError{}) {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": fmt.Sprintf("article with slug '%s' not found", slug),
						"error":   err.Error(),
					},
				})
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "get article by slug failed",
					"error":   err.Error(),
				},
			})
			return
		}

		deleted := false
		for i := range article.Comments {
			if article.Comments[i].ID == id {
				// check owner
				if article.Comments[i].AuthorID != currentUser.ID {
					w.Header().Set("Content-Type", "application/json; charset=utf-8")
					w.WriteHeader(http.StatusForbidden)
					_ = json.NewEncoder(w).Encode(ErrorResponse{
						Errors: map[string]interface{}{
							"message": "you are not author of this comment",
						},
					})
					return
				}

				article.Comments = append(article.Comments[:i], article.Comments[i+1:]...)
				deleted = true
				break
			}
		}

		if !deleted {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": fmt.Sprintf("comment with id '%d' in article with slug '%s' not found", id, slug),
				},
			})
			return
		}

		// update article
		err = h.articleRepo.UpdateBySlug(slug, *article)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "error on update article",
					"body":    err.Error(),
				},
			})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *handler) handleFavoriteArticle() http.HandlerFunc {
	type Response SingleArticleResponse

	return func(w http.ResponseWriter, r *http.Request) {
		// get current user
		currentUser := r.Context().Value(currentUserCtx).(*models.User)

		// get params
		slug := r.Context().Value("slug").(string)

		// get article by slug
		article, err := h.articleRepo.GetBySlug(slug)
		if err != nil {
			if errors.As(err, &models.UserByUsernameNotFoundError{}) {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": "article not found",
						"error":   err.Error(),
					},
				})
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "get article by slug failed",
					"error":   err.Error(),
				},
			})
			return
		}

		article.Favorites[currentUser.ID] = true

		// update article
		err = h.articleRepo.UpdateBySlug(slug, *article)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "update article failed",
					"error":   err.Error(),
				},
			})
			return
		}

		// find user by id
		user, err := h.userRepo.GetByID(article.AuthorID)
		if err != nil {
			if errors.As(err, &models.UserByIDNotFoundError{}) {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": "author of article not found",
						"error":   err.Error(),
					},
				})
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "get author of article failed",
					"error":   err.Error(),
				},
			})
			return
		}

		favorited := false
		if currentUser != nil {
			if f, ok := article.Favorites[currentUser.ID]; f && ok {
				favorited = true
			}
		}

		following := false
		if currentUser != nil {
			if f, ok := user.Followers[currentUser.ID]; f && ok {
				following = true
			}
		}

		// success response
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Response{
			Article: Article{
				Slug:           article.Slug,
				Title:          article.Title,
				Description:    article.Description,
				Body:           article.Body,
				TagList:        article.Tags,
				CreatedAt:      article.CreatedAt.UTC().Format(dateLayout),
				UpdatedAt:      article.UpdatedAt.UTC().Format(dateLayout),
				Favorited:      favorited,
				FavoritesCount: len(article.Favorites),
				Author: Author{
					Username:  user.Username,
					Bio:       user.Bio,
					Image:     user.Image,
					Following: following,
				},
			},
		})
	}
}

func (h *handler) handleUnfavoriteArticle() http.HandlerFunc {
	type Response SingleArticleResponse

	return func(w http.ResponseWriter, r *http.Request) {
		// get current user
		currentUser := r.Context().Value(currentUserCtx).(*models.User)

		// get params
		slug := r.Context().Value("slug").(string)

		// get article by slug
		article, err := h.articleRepo.GetBySlug(slug)
		if err != nil {
			if errors.As(err, &models.UserByUsernameNotFoundError{}) {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": "article not found",
						"error":   err.Error(),
					},
				})
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "get article by slug failed",
					"error":   err.Error(),
				},
			})
			return
		}

		delete(article.Favorites, currentUser.ID)

		// update article
		err = h.articleRepo.UpdateBySlug(slug, *article)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "update article failed",
					"error":   err.Error(),
				},
			})
			return
		}

		// find user by id
		user, err := h.userRepo.GetByID(article.AuthorID)
		if err != nil {
			if errors.As(err, &models.UserByIDNotFoundError{}) {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(ErrorResponse{
					Errors: map[string]interface{}{
						"message": "author of article not found",
						"error":   err.Error(),
					},
				})
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "get author of article failed",
					"error":   err.Error(),
				},
			})
			return
		}

		favorited := false
		if currentUser != nil {
			if f, ok := article.Favorites[currentUser.ID]; f && ok {
				favorited = true
			}
		}

		following := false
		if currentUser != nil {
			if f, ok := user.Followers[currentUser.ID]; f && ok {
				following = true
			}
		}

		// success response
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Response{
			Article: Article{
				Slug:           article.Slug,
				Title:          article.Title,
				Description:    article.Description,
				Body:           article.Body,
				TagList:        article.Tags,
				CreatedAt:      article.CreatedAt.UTC().Format(dateLayout),
				UpdatedAt:      article.UpdatedAt.UTC().Format(dateLayout),
				Favorited:      favorited,
				FavoritesCount: len(article.Favorites),
				Author: Author{
					Username:  user.Username,
					Bio:       user.Bio,
					Image:     user.Image,
					Following: following,
				},
			},
		})
	}
}

func (h *handler) handleGetTags() http.HandlerFunc {
	type Response ListOfTagsResponse

	return func(w http.ResponseWriter, r *http.Request) {
		// get tags
		tags, err := h.articleRepo.GetTags()
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Errors: map[string]interface{}{
					"message": "get tags failed",
					"error":   err.Error(),
				},
			})
			return
		}

		// success response
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Response{
			Tags: tags,
		})
	}
}
