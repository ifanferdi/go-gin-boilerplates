package bootstrap

import (
	"log"
	"strings"

	"go-gin-boilerplate/internal/config"
	"go-gin-boilerplate/internal/constants"
	"go-gin-boilerplate/internal/features/category"
	"go-gin-boilerplate/internal/features/product"
	grpcinfra "go-gin-boilerplate/internal/infrastructure/grpc"
	"go-gin-boilerplate/internal/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type App struct {
	Config     *config.Config
	HTTPServer *gin.Engine
	GRPCServer *grpcinfra.Server
	ExternalCC *grpc.ClientConn
}

func NewApp(cfg *config.Config, db *gorm.DB) (*App, error) {
	if strings.EqualFold(cfg.App.Env, "production") {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(middlewares.ErrorHandler(cfg))
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	router.GET("/", func(c *gin.Context) {
		c.JSON(constants.StatusOK, gin.H{
			"name":     cfg.App.Name,
			"version":  cfg.App.Version,
			"grpcPort": cfg.GRPC.Port,
			"endpoints": []string{
				"/api/v1/categories",
				"/api/v1/products",
			},
		})
	})

	grpcSrv := grpcinfra.NewServer()

	// Koneksi gRPC client di infrastructure (lazy connection via grpc.NewClient)
	var extConn *grpc.ClientConn
	if cfg.GRPC.ExternalServiceURL != "" {
		conn, err := grpcinfra.NewClientConnection(cfg.GRPC.ExternalServiceURL)
		if err != nil {
			log.Printf("⚠️  gRPC client connection warning: %v\n", err)
		} else {
			extConn = conn
		}
	}

	apiV1 := router.Group("/api/v1")
	{
		category.RegisterFeature(apiV1, grpcSrv.Engine, db)
		product.RegisterFeature(apiV1, grpcSrv.Engine, db)
	}

	router.NoRoute(middlewares.NotFoundHandler())

	return &App{
		Config:     cfg,
		HTTPServer: router,
		GRPCServer: grpcSrv,
		ExternalCC: extConn,
	}, nil
}

func (a *App) Close() {
	if a.GRPCServer != nil {
		a.GRPCServer.Stop()
	}
	if a.ExternalCC != nil {
		_ = a.ExternalCC.Close()
	}
}
