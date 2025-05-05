package rest

import (
	"github.com/Gen1usBruh/warehouse-api/internal/scope"

	_ "github.com/Gen1usBruh/warehouse-api/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type HandlerConfig struct {
	Dep *scope.Dependencies
}

func NewHandler(cfg HandlerConfig) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/products", cfg.CreateProduct)
	r.GET("/products/:id", cfg.GetProduct)
	r.PUT("/products/:id", cfg.UpdateProduct)
	r.DELETE("/products/:id", cfg.DeleteProduct)
	r.GET("/products", cfg.ListProducts)

	return r
}
