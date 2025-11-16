package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/SenechkaP/avito-test/configs"
	prHandler "github.com/SenechkaP/avito-test/internal/handler/pr"
	teamHandler "github.com/SenechkaP/avito-test/internal/handler/team"
	userHandler "github.com/SenechkaP/avito-test/internal/handler/users"

	prRepo "github.com/SenechkaP/avito-test/internal/repository/pr"
	teamRepo "github.com/SenechkaP/avito-test/internal/repository/team"
	userRepo "github.com/SenechkaP/avito-test/internal/repository/users"

	prService "github.com/SenechkaP/avito-test/internal/service/pr"
	teamService "github.com/SenechkaP/avito-test/internal/service/team"
	userService "github.com/SenechkaP/avito-test/internal/service/users"

	"github.com/SenechkaP/avito-test/internal/server"
	"github.com/go-chi/chi/v5"
)

func main() {
	log.SetOutput(os.Stdout)
	config := configs.LoadConfig(".env")

	router := chi.NewRouter()

	db := server.InitDB(&config.DbConfig)
	defer db.Close()

	prRepo := prRepo.NewPRRepository(db)
	prService := prService.NewPRService(prRepo)
	prHandler := prHandler.NewPRHandler(prService)
	prHandler.RegisterRoutes(router)

	teamRepo := teamRepo.NewTeamRepository(db)
	teamService := teamService.NewTeamService(teamRepo)
	teamHandler := teamHandler.NewTeamHandler(teamService)
	teamHandler.RegisterRoutes(router)

	userRepo := userRepo.NewUserRepository(db)
	userService := userService.NewUserService(userRepo)
	userHandler := userHandler.NewUserHandler(userService)
	userHandler.RegisterRoutes(router)

	addr := ":" + config.ServerConfig.Port
	srv := &http.Server{Addr: addr, Handler: router}

	go server.GracefulShutdown(srv)

	log.Printf("Server is launched on: %s\n", srv.Addr)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server start error: %v\n", err)
	}
}
