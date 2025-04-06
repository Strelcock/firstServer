package main

import (
	"bytes"
	"encoding/json"
	"firstServer/internal/auth"
	"firstServer/internal/user"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDb() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func initData(db *gorm.DB) {
	db.Create(&user.User{
		Email:    "a@a.ru",
		Password: "$2a$10$ICxhLrcQFa2l0WcXysl2nuKRButsk.kQsdNUo56znQBk4IF6G/8xG",
		Name:     "IVAn",
	})
}

func removeData(db *gorm.DB) {
	db.Unscoped().
		Where("email = ?", "a@a.ru").
		Delete(&user.User{})
}

func TestLoginSuccess(t *testing.T) {
	//Prep
	db := initDb()
	initData(db)

	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(auth.LoginRequest{
		Email:    "a@a.ru",
		Password: "qweasdzxc",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 200 {
		t.Fatalf("Expected %d got %d", 200, res.StatusCode)
	}

	var response auth.LoginResponse
	byteBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(byteBody, &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Token == "" {
		t.Fatalf("Token empty")
	}

	removeData(db)
}

func TestLoginFail(t *testing.T) {
	//Prep
	db := initDb()
	initData(db)

	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(auth.LoginRequest{
		Email:    "a@a.ru",
		Password: "qweasd1zxc",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode == 200 {
		t.Fatalf("Expected %d got %d", 401, res.StatusCode)
	}

	removeData(db)
}
