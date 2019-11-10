package inmem

import "github.com/nasermirzaei89/realworld-go/models"

type userRepo struct {
	users []models.User
}

func NewUserRepository() models.UserRepository {
	return &userRepo{
		users: make([]models.User, 0),
	}
}

func (repo *userRepo) NewID() int {
	panic("implement me")
}

func (repo *userRepo) GetByEmail(email string) (res *models.User, err error) {
	panic("implement me")
}

func (repo *userRepo) GetByUsername(username string) (res *models.User, err error) {
	panic("implement me")
}

func (repo *userRepo) GetByID(id int) (res *models.User, err error) {
	panic("implement me")
}

func (repo *userRepo) Add(user models.User) error {
	panic("implement me")
}

func (repo *userRepo) UpdateByID(id int, user models.User) (err error) {
	panic("implement me")
}

func (repo *userRepo) ListByFollowedBy(userID int) (res []models.User, err error) {
	panic("implement me")
}
