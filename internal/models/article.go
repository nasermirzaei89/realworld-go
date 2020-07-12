package models

import (
	"fmt"
	"time"
)

type Article struct {
	Slug        string // unique
	Title       string
	Description string
	Body        string
	Tags        []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	AuthorID    int
	Favorites   map[int]bool
	Comments    []Comment
}

type Comment struct {
	ID        int // unique
	CreatedAt time.Time
	UpdatedAt time.Time
	Body      string
	AuthorID  int
}

type ArticleRepository interface {
	List(offset, limit int, filters ...ArticleFilter) (res []Article, total int, err error)
	GetBySlug(slug string) (res *Article, err error)
	Add(entity Article) (err error)
	UpdateBySlug(slug string, entity Article) (err error)
	DeleteBySlug(slug string) (err error)
	NewCommentID() (id int)
	GetTags() (res []string, err error)
}

type ArticleBySlugNotFoundError struct {
	Slug string
}

func (e ArticleBySlugNotFoundError) Error() string {
	return fmt.Sprintf("article with slug '%s' not found", e.Slug)
}

type ArticleFilter func([]Article) []Article

func FilterArticlesByTag(tag string) ArticleFilter {
	return func(articles []Article) []Article {
		var res []Article
		for _, article := range articles {
			for i := range article.Tags {
				if tag == article.Tags[i] {
					res = append(res, article)
					break
				}
			}
		}

		return res
	}
}

func FilterArticlesByAuthor(user User) ArticleFilter {
	return func(articles []Article) []Article {
		var res []Article
		for _, article := range articles {
			if article.AuthorID == user.ID {
				res = append(res, article)
			}
		}

		return res
	}
}

func FilterArticlesByAuthors(users ...User) ArticleFilter {
	return func(articles []Article) []Article {
		var res []Article
		for _, article := range articles {
			for _, user := range users {
				if article.AuthorID == user.ID {
					res = append(res, article)
					break
				}
			}
		}

		return res
	}
}

func FilterArticlesByFavorite(user User) ArticleFilter {
	return func(articles []Article) []Article {
		var res []Article
		for _, article := range articles {
			if f, ok := article.Favorites[user.ID]; f && ok {
				res = append(res, article)
			}
		}

		return res
	}
}
