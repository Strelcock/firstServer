package auth

import (
	"firstServer/internal/di"
	"firstServer/internal/user"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo di.IUserRepository
}

func NewAuthService(userRepo di.IUserRepository) *AuthService {
	return &AuthService{
		UserRepo: userRepo,
	}
}

func (as *AuthService) Register(email, password, name string) (string, error) {
	existedUser, _ := as.UserRepo.FindByEmail(email)
	if existedUser != nil {
		return "", ErrUserExists
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user := &user.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPass),
	}

	_, err = as.UserRepo.Create(user)
	if err != nil {
		return "", err
	}
	return user.Email, nil
}

func (as *AuthService) Login(email, password string) (string, error) {
	existedUser, _ := as.UserRepo.FindByEmail(email)
	if existedUser == nil {
		return "", ErrWrongCredentials
	}

	err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(password))
	if err != nil {
		return "", ErrWrongCredentials
	}
	return email, nil
}
