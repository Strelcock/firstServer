package main

import (
	"firstServer/configs"
	"firstServer/internal/auth"
	"firstServer/internal/link"
	"firstServer/internal/user"
	"firstServer/pkg/db"
	"firstServer/pkg/middleware"
	"fmt"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()

	db := db.NewDb(conf)
	router := http.NewServeMux()

	//repos
	linkRepo := link.NewLinkRepository(db)
	userRepo := user.NewUserRepository(db)

	//services
	authService := auth.NewAuthService(userRepo)

	//handler
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})

	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepo: linkRepo,
		Config:   conf,
	})

	//Midlewares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8081",
		Handler: stack(router),
	}

	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()

}
