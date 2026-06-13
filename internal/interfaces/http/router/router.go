package router

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	v1 "github.com/infracore/infracore/internal/interfaces/http/handler/v1"
	"github.com/infracore/infracore/internal/interfaces/http/middleware"
	"github.com/infracore/infracore/pkg/config"
	"github.com/infracore/infracore/pkg/crypto"
	"go.uber.org/zap"
)

// Deps bundles all handler dependencies for the router.
type Deps struct {
	AuthHandler       *v1.AuthHandler
	MonitoringHandler *v1.MonitoringHandler
	JWTManager        *crypto.JWTManager
	Cfg               *config.Config
	Log               *zap.Logger
}

// New assembles and returns the Gin engine with all routes registered.
func New(d Deps) *gin.Engine {
	if d.Cfg.IsDevelopment() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// ── Global middleware ────────────────────────────────────────
	r.Use(middleware.Recovery(d.Log))
	r.Use(middleware.RequestLogger(d.Log))
	r.Use(cors.New(buildCORSConfig(d.Cfg.CORS)))
	r.MaxMultipartMemory = int64(d.Cfg.Server.MaxBodyMB) << 20

	// ── Health ──────────────────────────────────────────────────
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "version": d.Cfg.App.Version})
	})
	r.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	// ── API v1 ──────────────────────────────────────────────────
	api := r.Group("/api/v1")

	// Public auth routes (no JWT required)
	auth := api.Group("/auth")
	{
		auth.POST("/login", d.AuthHandler.Login)
		auth.POST("/refresh", d.AuthHandler.Refresh)
		auth.POST("/logout", d.AuthHandler.Logout)
	}

	// Protected routes (JWT required)
	protected := api.Group("")
	protected.Use(middleware.Auth(d.JWTManager))
	{
		// Auth self-service
		protected.GET("/auth/me", d.AuthHandler.Me)

		// ── DCIM ────────────────────────────────────────────────
		dcim := protected.Group("/dcim")
		_ = dcim // handlers registered in Phase 3
		/*
			dcim.GET("/data-centers",     RequirePermission("dcim.datacenters","list"),   d.DCIMHandler.ListDataCenters)
			dcim.POST("/data-centers",    RequirePermission("dcim.datacenters","create"),  d.DCIMHandler.CreateDataCenter)
			dcim.GET("/data-centers/:id", RequirePermission("dcim.datacenters","read"),    d.DCIMHandler.GetDataCenter)
			dcim.PUT("/data-centers/:id", RequirePermission("dcim.datacenters","update"),  d.DCIMHandler.UpdateDataCenter)
			dcim.DELETE("/data-centers/:id", RequirePermission("dcim.datacenters","delete"), d.DCIMHandler.DeleteDataCenter)
			dcim.GET("/racks",            RequirePermission("dcim.racks","list"),          d.DCIMHandler.ListRacks)
			dcim.POST("/racks",           RequirePermission("dcim.racks","create"),        d.DCIMHandler.CreateRack)
		*/

		// ── Assets ──────────────────────────────────────────────
		assets := protected.Group("/assets")
		_ = assets
		/*
			assets.GET("",     RequirePermission("asset.assets","list"),   d.AssetHandler.List)
			assets.POST("",    RequirePermission("asset.assets","create"),  d.AssetHandler.Create)
			assets.GET("/:id", RequirePermission("asset.assets","read"),    d.AssetHandler.Get)
			assets.PUT("/:id", RequirePermission("asset.assets","update"),  d.AssetHandler.Update)
		*/

		// ── IPAM ────────────────────────────────────────────────
		ipam := protected.Group("/ipam")
		_ = ipam

		// ── Monitoring ──────────────────────────────────────────
		mon := protected.Group("/monitoring")
		if d.MonitoringHandler != nil {
			mon.GET("/hosts/counts", d.MonitoringHandler.GetStatusCounts)
			mon.GET("/hosts", d.MonitoringHandler.ListHosts)
			mon.POST("/hosts", d.MonitoringHandler.CreateHost)
			mon.GET("/hosts/:id", d.MonitoringHandler.GetHost)
			mon.DELETE("/hosts/:id", d.MonitoringHandler.DeleteHost)
		} else {
			_ = mon
		}

		// ── Alerting ────────────────────────────────────────────
		alerting := protected.Group("/alerting")
		_ = alerting

		// ── Agents ──────────────────────────────────────────────
		agentGroup := protected.Group("/agents")
		_ = agentGroup

		// ── Licenses ────────────────────────────────────────────
		licenses := protected.Group("/licenses")
		_ = licenses

		// ── Contracts ───────────────────────────────────────────
		contracts := protected.Group("/contracts")
		_ = contracts

		// ── Discovery ───────────────────────────────────────────
		discovery := protected.Group("/discovery")
		_ = discovery

		// ── IAM ─────────────────────────────────────────────────
		iam := protected.Group("/iam")
		_ = iam
	}

	// 404 fallback
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "route not found"})
	})

	return r
}

func buildCORSConfig(cfg config.CORSConfig) cors.Config {
	c := cors.DefaultConfig()
	c.AllowOrigins = cfg.AllowedOrigins
	c.AllowMethods = cfg.AllowedMethods
	c.AllowHeaders = cfg.AllowedHeaders
	c.ExposeHeaders = cfg.ExposeHeaders
	c.AllowCredentials = true // needed for httpOnly refresh cookie
	// The YAML value is expressed in seconds; gin-contrib/cors requires a duration.
	c.MaxAge = time.Duration(cfg.MaxAge) * time.Second
	return c
}
