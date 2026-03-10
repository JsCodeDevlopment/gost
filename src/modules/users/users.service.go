package users

import (
	"gost/src/modules/users/dto"
	"gost/src/modules/users/entities"
)

type UserService interface {
	Create(d dto.CreateUserDto) (*entities.User, error)
	FindAll() ([]entities.User, error)
	FindById(id uint) (*entities.User, error)
	Update(id uint, user *entities.User) error
	Delete(id uint) error
	UpdateAvatar(id uint, avatarPath string) error
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) Create(d dto.CreateUserDto) (*entities.User, error) {
	user := &entities.User{
		Name:  d.Name,
		Email: d.Email,
	}
	err := s.repo.Create(user)
	return user, err
}

func (s *userService) FindAll() ([]entities.User, error) {
	return s.repo.FindAll()
}

func (s *userService) FindById(id uint) (*entities.User, error) {
	return s.repo.FindById(id)
}

func (s *userService) Update(id uint, user *entities.User) error {
	existing, err := s.FindById(id)
	if err != nil {
		return err
	}
	// Copy updatable fields
	existing.Name = user.Name
	existing.Email = user.Email
	return s.repo.Update(existing)
}

func (s *userService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *userService) UpdateAvatar(id uint, avatarPath string) error {
	existing, err := s.FindById(id)
	if err != nil {
		return err
	}
	existing.Avatar = avatarPath
	return s.repo.Update(existing)
}
