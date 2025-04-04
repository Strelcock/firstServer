package main

import (
	"firstServer/configs"
	"firstServer/internal/auth"
	"firstServer/internal/link"
	"firstServer/internal/stat"
	"firstServer/internal/user"
	"firstServer/pkg/db"
	"firstServer/pkg/event"
	"firstServer/pkg/middleware"
	"fmt"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()

	db := db.NewDb(conf)
	router := http.NewServeMux()
	eventBus := event.NewEventBus()

	//repos
	linkRepo := link.NewLinkRepository(db)
	userRepo := user.NewUserRepository(db)
	statRepo := stat.NewStatRepository(db)

	//services
	authService := auth.NewAuthService(userRepo)
	statService := stat.NewStatService(stat.StatServiceDeps{
		EventBus:       eventBus,
		StatRepository: statRepo,
	})

	//handler
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})

	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepo: linkRepo,
		EventBus: eventBus,
		Config:   conf,
	})

	stat.NewStatHandler(router, stat.StatHandlerDeps{
		StatRepo: statRepo,
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

	go statService.AddClick()

	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()

}
