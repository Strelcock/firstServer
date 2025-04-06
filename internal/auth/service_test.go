package auth_test

import (
	"firstServer/internal/auth"
	"firstServer/internal/user"
	"testing"
)

type MockUserRepository struct{}

func (mur *MockUserRepository) Create(u *user.User) (*user.User, error) {
	return &user.User{
		Email: "a@a.ru",
	}, nil
}

func (mur *MockUserRepository) FindByEmail(email string) (*user.User, error) {
	return nil, nil
}

func TestRegisterSuccess(t *testing.T) {
	const initialEmail = "a@a.ru"
	authService := auth.NewAuthService(&MockUserRepository{})
	email, err := authService.Register(initialEmail, "1", "ВАСЯ")
	if err != nil {
		t.Fatal(err)
	}

	if email != initialEmail {
		t.Fatalf("Email %s do not match %s", email, initialEmail)
	}

}
