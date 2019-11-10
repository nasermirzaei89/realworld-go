package inmem

import "github.com/nasermirzaei89/realworld-go/models"

type articleRepo struct {
	articles []models.Article
}

func NewArticleRepository() models.ArticleRepository {
	return &articleRepo{
		articles: make([]models.Article, 0),
	}
}

func (repo *articleRepo) List(offset, limit int, filters ...models.ArticleFilter) (res []models.Article, total int, err error) {
	panic("implement me")
}

func (repo *articleRepo) GetBySlug(slug string) (res *models.Article, err error) {
	panic("implement me")
}

func (repo *articleRepo) Create(article models.Article) (err error) {
	panic("implement me")
}

func (repo *articleRepo) UpdateBySlug(slug string, article models.Article) (err error) {
	panic("implement me")
}

func (repo *articleRepo) DeleteBySlug(slug string) (err error) {
	panic("implement me")
}

func (repo *articleRepo) NewCommentID() (id int) {
	panic("implement me")
}

func (repo *articleRepo) GetTags() (res []string, err error) {
	panic("implement me")
}
