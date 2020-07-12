package inmem

import (
	"fmt"
	"github.com/nasermirzaei89/realworld-go/internal/models"
)

type userRepo struct {
	users  []models.User
	nextID int
}

func NewUserRepository() models.UserRepository {
	return &userRepo{
		users:  make([]models.User, 0),
		nextID: 1,
	}
}

func (repo *userRepo) NewID() int {
	defer func() { repo.nextID = repo.nextID + 1 }()
	return repo.nextID
}

func (repo *userRepo) GetByEmail(email string) (*models.User, error) {
	for _, user := range repo.users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, models.UserByEmailNotFoundError{Email: email}
}

func (repo *userRepo) GetByUsername(username string) (*models.User, error) {
	for _, user := range repo.users {
		if user.Username == username {
			return &user, nil
		}
	}

	return nil, models.UserByUsernameNotFoundError{Username: username}
}

func (repo *userRepo) GetByID(id int) (*models.User, error) {
	for _, user := range repo.users {
		if user.ID == id {
			return &user, nil
		}
	}

	return nil, models.UserByIDNotFoundError{ID: id}
}

func (repo *userRepo) Add(entity models.User) error {
	for _, user := range repo.users {
		if user.ID == entity.ID {
			return fmt.Errorf("user with id '%d' already exists", entity.ID)
		}
		if user.Email == entity.Email {
			return fmt.Errorf("user with email '%s' already exists", entity.Email)
		}
		if user.Username == entity.Username {
			return fmt.Errorf("user with username '%s' already exists", entity.Username)
		}
		if entity.Token != "" && user.Token == entity.Token {
			return fmt.Errorf("user with token '%s' already exists", entity.Token)
		}
	}

	repo.users = append(repo.users, entity)

	return nil
}

func (repo *userRepo) UpdateByID(id int, entity models.User) error {
	index := -1
	for i, user := range repo.users {
		if user.ID == id {
			index = i
			continue
		}
		if user.ID == entity.ID {
			return fmt.Errorf("user with id '%d' already exists", entity.ID)
		}
		if user.Email == entity.Email {
			return fmt.Errorf("user with email '%s' already exists", entity.Email)
		}
		if user.Username == entity.Username {
			return fmt.Errorf("user with username '%s' already exists", entity.Username)
		}
		if entity.Token != "" && user.Token == entity.Token {
			return fmt.Errorf("user with token '%s' already exists", entity.Token)
		}
	}

	if index == -1 {
		return models.UserByIDNotFoundError{ID: id}
	}

	repo.users[index] = entity

	return nil
}

func (repo *userRepo) ListByFollowedBy(userID int) ([]models.User, error) {
	var res []models.User
	for _, user := range repo.users {
		if f, ok := user.Followers[userID]; f && ok {
			res = append(res, user)
		}
	}

	return res, nil
}
