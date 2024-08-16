package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mirhijinam/backend-bootcamp-assignment-2024/generated"
	"github.com/mirhijinam/backend-bootcamp-assignment-2024/internal/config"
	"github.com/mirhijinam/backend-bootcamp-assignment-2024/internal/net"
	"github.com/mirhijinam/backend-bootcamp-assignment-2024/internal/net/middleware"
	"github.com/mirhijinam/backend-bootcamp-assignment-2024/internal/pkg/db"
	"github.com/mirhijinam/backend-bootcamp-assignment-2024/internal/pkg/logger"
	authrepo "github.com/mirhijinam/backend-bootcamp-assignment-2024/internal/repository/auth"
	flatsrepo "github.com/mirhijinam/backend-bootcamp-assignment-2024/internal/repository/flats"
	houserepo "github.com/mirhijinam/backend-bootcamp-assignment-2024/internal/repository/houses"
	authservice "github.com/mirhijinam/backend-bootcamp-assignment-2024/internal/service/auth"
	flatsservice "github.com/mirhijinam/backend-bootcamp-assignment-2024/internal/service/flats"
	housesservice "github.com/mirhijinam/backend-bootcamp-assignment-2024/internal/service/houses"
	"go.uber.org/zap"
)

func Run() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	lgr := logger.New(cfg.LoggerConfig.Mode, cfg.LoggerConfig.Filepath)
	defer lgr.Sync()

	pool, err := db.MustOpenDB(context.Background(), cfg.DBConfig)
	if err != nil {
		panic(err)
	}

	authService := authservice.New(
		authrepo.New(pool, lgr), lgr)

	houseService := housesservice.New(
		houserepo.New(pool, lgr), lgr)

	flatsService := flatsservice.New(
		flatsrepo.New(pool, lgr), lgr)

	authMiddl := middleware.New(lgr)

	server, err := generated.NewServer(
		net.NewHandler(authService, houseService, flatsService, lgr),
		authMiddl,
	)
	if err != nil {
		panic(err)
	}

	lgr.Info("Server is running", zap.String("port", cfg.ServerConfig.Port))
	if err := http.ListenAndServe(
		fmt.Sprintf(":%s", cfg.ServerConfig.Port),
		server,
	); err != nil {
		panic(err)
	}
}
