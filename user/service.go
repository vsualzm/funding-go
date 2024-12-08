package user

import "golang.org/x/crypto/bcrypt"

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
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
