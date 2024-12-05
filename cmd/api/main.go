package main

import (
	"context"
	"log"
	"medods-test/internal/auth/repo/postgres"
	"medods-test/internal/auth/rest"
	"medods-test/internal/auth/service"
	"medods-test/internal/config"
	"medods-test/pkg/auth"
	"medods-test/pkg/db"
	"net/http"
)

func main() {
	ctx := context.Background()
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	DB, err := db.OpenDB(ctx, cfg.DBConfig)
	if err != nil {
		log.Fatalln(err)
	}

	userRepo := postgres.NewUserRepo(DB)
	sessionRepo := postgres.NewSessionRepo(DB)

	repo := &service.Repository{
		UserRepo:    userRepo,
		SessionRepo: sessionRepo,
	}

	s := service.New(repo)

	manager, err := auth.NewManager(cfg.AuthConfig.SigningKey)
	if err != nil {
		log.Fatalln(err)
	}

	restUseCase := &rest.UseCase{
		User: s.User(
			manager,
			cfg.AuthConfig.AccessTokenTTL,
			cfg.AuthConfig.RefreshTokenTTL,
		),
	}

	h := rest.New(restUseCase)

	server := &http.Server{
		Addr:    cfg.ServerConfig.Address(),
		Handler: h.Handler(),
	}
	log.Printf("listening on %v", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}

}
