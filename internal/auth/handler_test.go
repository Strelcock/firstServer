package auth_test

import (
	"bytes"
	"encoding/json"
	"firstServer/configs"
	"firstServer/internal/auth"
	"firstServer/internal/user"
	"firstServer/pkg/db"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func bootstrap() (*auth.AuthHandler, sqlmock.Sqlmock, error) {
	database, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: database,
	}))
	if err != nil {
		return nil, nil, err
	}

	userRepo := user.NewUserRepository(&db.Db{
		DB: gormDb,
	})

	handler := auth.AuthHandler{
		Config: &configs.Config{
			Auth: configs.AuthConfig{
				Secret: "secret",
			},
		},
		AuthService: auth.NewAuthService(userRepo),
	}
	return &handler, mock, nil
}

func TestLoginSuccess(t *testing.T) {
	handler, mock, err := bootstrap()
	rows := sqlmock.NewRows([]string{"email", "password"}).
		AddRow("a@a.ru", "$2a$10$ICxhLrcQFa2l0WcXysl2nuKRButsk.kQsdNUo56znQBk4IF6G/8xG")

	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	if err != nil {
		t.Fatal(err)
	}

	data, _ := json.Marshal(auth.LoginRequest{
		Email:    "a@a.ru",
		Password: "qweasdzxc",
	})
	reader := bytes.NewReader(data)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", reader)
	handler.Login()(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Got %d expected %d", w.Code, 201)
	}
}

func TestRegisterHandlerSuccess(t *testing.T) {
	handler, mock, err := bootstrap()
	if err != nil {
		t.Fatal(err)
	}

	rows := sqlmock.NewRows([]string{"email", "password", "name"})
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	data, _ := json.Marshal(auth.RegisterRequest{
		Name:     "IVAn",
		Email:    "a@a.ru",
		Password: "qweasdzxc",
	})
	reader := bytes.NewReader(data)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/register", reader)
	handler.Register()(w, req)
	if w.Code != http.StatusAccepted {
		t.Fatalf("Got %d expected %d", w.Code, 201)
	}
}
