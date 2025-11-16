package server

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/SenechkaP/avito-test/configs"
	"github.com/SenechkaP/avito-test/internal/database"
)

func GracefulShutdown(srv *http.Server) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	<-ctx.Done()

	log.Println("Shutting down service-courier")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown error: %v\n", err)
	}
}

func InitDB(cfg *configs.DbConfig) *database.DB {
	dbCtx, dbCancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer dbCancel()

	db, err := database.NewDB(dbCtx, cfg)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	return db
}
