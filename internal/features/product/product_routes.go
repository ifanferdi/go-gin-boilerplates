package product

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func RegisterFeature(rg *gin.RouterGroup, grpcServer *grpc.Server, db *gorm.DB) {
	repo := NewProductRepository(db)
	svc := NewProductService(repo)

	RegisterHTTPRoutes(rg, svc)
	if grpcServer != nil {
		RegisterGRPCRoutes(grpcServer, svc)
	}
}

func RegisterHTTPRoutes(rg *gin.RouterGroup, svc ProductService) {
	ctrl := NewProductController(svc)

	products := rg.Group("/products")
	{
		products.POST("/get", ctrl.FindAll)
		products.GET("", ctrl.FindAll)
		products.GET("/:id", ctrl.FindOne)
		products.POST("", ctrl.Create)
		products.PUT("/:id", ctrl.Update)
		products.DELETE("/:id", ctrl.Destroy)
	}
}

func RegisterGRPCRoutes(grpcServer *grpc.Server, svc ProductService) {
	grpcHandler := NewProductGRPCServer(svc)
	grpcServer.RegisterService(&ProductServiceDesc, grpcHandler)
}
