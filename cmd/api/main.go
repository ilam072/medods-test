package main

import (
	"context"
	"medods-test/internal/auth/repo/postgres"
	"medods-test/internal/auth/rest"
	"medods-test/internal/auth/service"
	"medods-test/internal/config"
	"medods-test/pkg/auth"
	"medods-test/pkg/db"
	"medods-test/pkg/email/smtp"
	"medods-test/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Error(err)
		return
	}

	DB, err := db.OpenDB(ctx, cfg.DBConfig)
	if err != nil {
		logger.Error(err)
		return
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
		logger.Error(err)
		return
	}

	smtpSender, err := smtp.NewSMTPSender(
		cfg.SMTPConfig.From,
		cfg.SMTPConfig.Pass,
		cfg.SMTPConfig.Host,
		cfg.SMTPConfig.Port,
	)
	if err != nil {
		logger.Error(err)
		return
	}

	restUseCase := &rest.UseCase{
		User: s.User(
			manager,
			smtpSender,
			cfg.AuthConfig.AccessTokenTTL,
			cfg.AuthConfig.RefreshTokenTTL,
		),
	}

	h := rest.New(restUseCase)

	server := &http.Server{
		Addr:    cfg.ServerConfig.Address(),
		Handler: h.Handler(),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Error("failed to run http server: %s\n", err.Error())
		}
	}()

	logger.Info("server started")

	waitForShutdown(server)
}

func waitForShutdown(server *http.Server) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	logger.Info("shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("failed to shutdown server: %v", err)
	}

	select {
	case <-ctx.Done():
		logger.Info("Done")
	}
}
