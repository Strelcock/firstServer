package main

import (
	"context"
	"firstServer/configs"
	"firstServer/internal/auth"
	"firstServer/internal/link"
	"firstServer/internal/user"
	"firstServer/pkg/db"
	"firstServer/pkg/middleware"
	"fmt"
	"net/http"
	"time"
)

func tickOperation(ctx context.Context) {
	ticker := time.NewTicker(200 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			fmt.Println("Tick")
		case <-ctx.Done():
			fmt.Println("Cancel")
			return
		}
	}
}

func main2() {
	ctx, cancel := context.WithCancel(context.Background())

	go tickOperation(ctx)
	time.Sleep(2 * time.Second)
	cancel()
	time.Sleep(2 * time.Second)
}

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
