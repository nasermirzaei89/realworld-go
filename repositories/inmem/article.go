package inmem

import (
	"fmt"
	"github.com/nasermirzaei89/realworld-go/models"
)

type articleRepo struct {
	articles      []models.Article
	nextCommentID int
}

func NewArticleRepository() models.ArticleRepository {
	return &articleRepo{
		articles:      make([]models.Article, 0),
		nextCommentID: 1,
	}
}

func (repo *articleRepo) List(offset, limit int, filters ...models.ArticleFilter) ([]models.Article, int, error) {
	res := make([]models.Article, len(repo.articles))
	copy(res, repo.articles)
	for _, filter := range filters {
		res = filter(res)
	}

	total := len(res)

	if offset > total {
		offset = total
	}

	if offset+limit > total {
		limit = total - offset
	}

	return res[offset : offset+limit], total, nil
}

func (repo *articleRepo) GetBySlug(slug string) (*models.Article, error) {
	for _, article := range repo.articles {
		if article.Slug == slug {
			return &article, nil
		}
	}

	return nil, &models.ArticleBySlugNotFoundError{Slug: slug}
}

func (repo *articleRepo) Add(entity models.Article) error {
	for _, article := range repo.articles {
		if article.Slug == entity.Slug {
			return fmt.Errorf("article with slug '%s' already exists", entity.Slug)
		}
	}

	repo.articles = append(repo.articles, entity)

	return nil
}

func (repo *articleRepo) UpdateBySlug(slug string, entity models.Article) error {
	index := -1
	for i, article := range repo.articles {
		if article.Slug == slug {
			index = i
			continue
		}
		if article.Slug == entity.Slug {
			return fmt.Errorf("article with slug '%s' already exists", entity.Slug)
		}
	}

	if index == -1 {
		return &models.ArticleBySlugNotFoundError{Slug: slug}
	}

	repo.articles[index] = entity

	return nil
}

func (repo *articleRepo) DeleteBySlug(slug string) error {
	for i, article := range repo.articles {
		if article.Slug == slug {
			repo.articles = append(repo.articles[:i], repo.articles[i+1:]...)
			return nil
		}
	}

	return &models.ArticleBySlugNotFoundError{Slug: slug}
}

func (repo *articleRepo) NewCommentID() int {
	defer func() { repo.nextCommentID = repo.nextCommentID + 1 }()
	return repo.nextCommentID
}

func (repo *articleRepo) GetTags() ([]string, error) {
	res := make([]string, 0)
	keys := make(map[string]bool)
	for _, article := range repo.articles {
		for _, tag := range article.Tags {
			if _, exists := keys[tag]; !exists {
				keys[tag] = true
				res = append(res, tag)
			}
		}
	}

	return res, nil
}
