package category

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func RegisterFeature(rg *gin.RouterGroup, grpcServer *grpc.Server, db *gorm.DB) CategoryRepository {
	repo := NewCategoryRepository(db)
	svc := NewCategoryService(repo)

	RegisterHTTPRoutes(rg, svc)
	if grpcServer != nil {
		RegisterGRPCRoutes(grpcServer, svc)
	}

	return repo
}

func RegisterHTTPRoutes(rg *gin.RouterGroup, svc CategoryService) {
	ctrl := NewCategoryController(svc)

	categories := rg.Group("/categories")
	{
		categories.POST("/get", ctrl.FindAll)
		categories.GET("", ctrl.FindAll)
		categories.GET("/:id", ctrl.FindOne)
		categories.POST("", ctrl.Create)
		categories.PUT("/:id", ctrl.Update)
		categories.DELETE("/:id", ctrl.Destroy)
	}
}

func RegisterGRPCRoutes(grpcServer *grpc.Server, svc CategoryService) {
	grpcHandler := NewCategoryGRPCServer(svc)
	grpcServer.RegisterService(&CategoryServiceDesc, grpcHandler)
}
