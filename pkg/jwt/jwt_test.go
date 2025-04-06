package jwt_test

import (
	"firstServer/pkg/jwt"
	"testing"
)

func TestJWTCreate(t *testing.T) {
	const email = "a@a.ru"
	jwtService := jwt.NewJWT("u/ZeBjl/3tIyLxkLJ0Xx9HQD27hpZnVc5IT+TKUHCyU=")
	token, err := jwtService.Create(jwt.JWTData{
		Email: email,
	})
	if err != nil {
		t.Fatal(err)
	}

	isValid, data := jwtService.Parse(token)
	if !isValid {
		t.Fatal("token is invalid")
	}

	if data.Email != email {
		t.Fatalf("Email %s not equal to %s", data.Email, email)
	}
}
