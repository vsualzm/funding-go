package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginInput) (User, error)
	CheckEmailAvailability(input CheckEmailInput) (bool, error)
	SaveAvatar(ID int, fileLocation string) (User, error)
	GetUserByID(ID int) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	user := User{}

	user.Name = input.Name
	user.Occupation = input.Occupation
	user.Email = input.Email
	// password hash to bycrypt
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	// sesudah di hash
	user.PasswordHash = string(passwordHash)
	// hardcode
	user.Role = "user"

	// save
	newUser, err := s.repository.Save(user)

	// return nil
	if err != nil {
		return newUser, err
	}

	// return service
	return newUser, nil

}

func (s *service) Login(input LoginInput) (User, error) {

	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return user, err
	}

	// cheking id user
	if user.ID == 0 {
		return user, errors.New("not user found that email")
	}

	// cheking password comparehash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *service) CheckEmailAvailability(input CheckEmailInput) (bool, error) {

	email := input.Email

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return false, err
	}

	if user.ID == 0 {
		return true, nil
	}

	return false, nil
}

func (s *service) SaveAvatar(ID int, fileLocation string) (User, error) {

	user, err := s.repository.FindById(ID)

	if err != nil {
		return user, err
	}

	user.AvatarFileName = fileLocation

	updateUser, err := s.repository.Update(user)

	if err != nil {
		return updateUser, err
	}

	return updateUser, nil
}

func (s *service) GetUserByID(ID int) (User, error) {

	user, err := s.repository.FindById(ID)

	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("user not found")
	}

	return user, nil

}
