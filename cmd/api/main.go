package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	appidentity "github.com/infracore/infracore/internal/application/identity"
	appmonitoring "github.com/infracore/infracore/internal/application/monitoring"
	"github.com/infracore/infracore/internal/infrastructure/cache/redis"
	"github.com/infracore/infracore/internal/infrastructure/persistence/postgres"
	pgidentity "github.com/infracore/infracore/internal/infrastructure/persistence/postgres/identity"
	pgmonitoring "github.com/infracore/infracore/internal/infrastructure/persistence/postgres/monitoring"
	v1 "github.com/infracore/infracore/internal/interfaces/http/handler/v1"
	"github.com/infracore/infracore/internal/interfaces/http/router"
	"github.com/infracore/infracore/pkg/config"
	"github.com/infracore/infracore/pkg/crypto"
	"github.com/infracore/infracore/pkg/logger"
)

func main() {
	// ── Config ──────────────────────────────────────────────────
	cfgFile := os.Getenv("INFRACORE_CONFIG")
	cfg, err := config.Load(cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	// ── Logger ──────────────────────────────────────────────────
	log, err := logger.New(cfg.Logging.Level, cfg.Logging.Format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "init logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync() //nolint:errcheck

	log.Info("starting InfraCore API",
		zap.String("version", cfg.App.Version),
		zap.String("env", cfg.App.Environment),
		zap.String("addr", cfg.Server.Addr()),
	)

	// ── Database ─────────────────────────────────────────────────
	db, err := postgres.NewDB(cfg.Database)
	if err != nil {
		log.Fatal("connect postgres", zap.Error(err))
	}
	defer db.Close()
	log.Info("postgres connected")

	// ── Redis ────────────────────────────────────────────────────
	_, err = redis.NewClient(cfg.Redis)
	if err != nil {
		log.Fatal("connect redis", zap.Error(err))
	}
	log.Info("redis connected")

	// ── Repositories ─────────────────────────────────────────────
	tenantRepo  := pgidentity.NewTenantRepository(db)
	userRepo    := pgidentity.NewUserRepository(db)
	roleRepo    := pgidentity.NewRoleRepository(db)
	tokenRepo   := pgidentity.NewRefreshTokenRepository(db)
	hostRepo    := pgmonitoring.NewHostRepository(db)

	// ── JWT manager ──────────────────────────────────────────────
	jwtMgr := crypto.NewJWTManager(
		cfg.Auth.JWTAccessSecret,
		cfg.Auth.JWTRefreshSecret,
		cfg.Auth.AccessTokenTTL,
		cfg.Auth.RefreshTokenTTL,
	)

	// ── Application services ─────────────────────────────────────
	authSvc       := appidentity.NewAuthService(tenantRepo, userRepo, roleRepo, tokenRepo, jwtMgr, cfg.Auth, log, db)
	monitoringSvc := appmonitoring.NewMonitoringService(hostRepo, log)

	// ── Handlers ─────────────────────────────────────────────────
	authHandler       := v1.NewAuthHandler(authSvc, cfg.Auth)
	monitoringHandler := v1.NewMonitoringHandler(monitoringSvc)

	// ── Router ───────────────────────────────────────────────────
	engine := router.New(router.Deps{
		AuthHandler:       authHandler,
		MonitoringHandler: monitoringHandler,
		JWTManager:        jwtMgr,
		Cfg:               cfg,
		Log:               log,
	})

	// ── HTTP server ───────────────────────────────────────────────
	srv := &http.Server{
		Addr:         cfg.Server.Addr(),
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		log.Info("API server listening", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen error", zap.Error(err))
		}
	}()

	// ── Graceful shutdown ─────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("shutdown error", zap.Error(err))
	}
	log.Info("server stopped")
}
