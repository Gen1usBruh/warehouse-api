package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Gen1usBruh/warehouse-api/internal/logger/sl"
	db "github.com/Gen1usBruh/warehouse-api/internal/storage/postgres/sqlc"
	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	Success   bool   `json:"success"`
	Error     string `json:"error,omitempty"`
	ErrorCode int    `json:"errorCode"`
}

type ProductRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=255"`
	Description string `json:"description" binding:"required,max=1000"`
	Price       int32  `json:"price" binding:"required,gt=0"`
	Quantity    int32  `json:"quantity" binding:"required,gte=0"`
}

func (h *HandlerConfig) CreateProduct(c *gin.Context) {
	const op = "internal.rest.product.create"
	var req ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Dep.Sl.Error(fmt.Sprintf("%s | Error binding: ", op), sl.Err(err))
		c.JSON(http.StatusBadRequest, BaseResponse{Error: "Wrong data", ErrorCode: 400})
		return
	}
	arg := db.CreateProductParams{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Quantity:    req.Quantity,
	}
	id, err := h.Dep.Db.CreateProduct(c, arg)
	if err != nil {
		h.Dep.Sl.Error(fmt.Sprintf("%s | Error creating product: ", op), sl.Err(err))
		c.JSON(http.StatusInternalServerError, BaseResponse{Error: "Failed to create product", ErrorCode: 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *HandlerConfig) GetProduct(c *gin.Context) {
	const op = "internal.rest.product.get"

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.Dep.Sl.Error(fmt.Sprintf("%s | Invalid ID param: ", op), sl.Err(err))
		c.JSON(http.StatusBadRequest, BaseResponse{Error: "Invalid product ID", ErrorCode: 400})
		return
	}

	product, err := h.Dep.Db.GetProductByID(c, int32(id))
	if err != nil {
		h.Dep.Sl.Error(fmt.Sprintf("%s | Product not found: ", op), sl.Err(err))
		c.JSON(http.StatusNotFound, BaseResponse{Error: "Product not found", ErrorCode: 404})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func (h *HandlerConfig) UpdateProduct(c *gin.Context) {
	const op = "internal.rest.product.update"

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.Dep.Sl.Error(fmt.Sprintf("%s | Invalid ID param: ", op), sl.Err(err))
		c.JSON(http.StatusBadRequest, BaseResponse{Error: "Invalid product ID", ErrorCode: 400})
		return
	}
	var req ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Dep.Sl.Error(fmt.Sprintf("%s | Error binding: ", op), sl.Err(err))
		c.JSON(http.StatusBadRequest, BaseResponse{Error: "Wrong data", ErrorCode: 400})
		return
	}

	arg := db.UpdateProductParams{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Quantity:    req.Quantity,
	}
	arg.ID = int32(id)

	if err := h.Dep.Db.UpdateProduct(c, arg); err != nil {
		h.Dep.Sl.Error(fmt.Sprintf("%s | Failed to update product: ", op), sl.Err(err))
		c.JSON(http.StatusInternalServerError, BaseResponse{Error: "Failed to update product", ErrorCode: 500})
		return
	}

	c.JSON(http.StatusOK, BaseResponse{Success: true})
}

func (h *HandlerConfig) DeleteProduct(c *gin.Context) {
	const op = "internal.rest.product.delete"

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.Dep.Sl.Error(fmt.Sprintf("%s | Invalid ID param: ", op), sl.Err(err))
		c.JSON(http.StatusBadRequest, BaseResponse{Error: "Invalid product ID", ErrorCode: 400})
		return
	}

	if err := h.Dep.Db.DeleteProduct(c, int32(id)); err != nil {
		h.Dep.Sl.Error(fmt.Sprintf("%s | Failed to delete product: ", op), sl.Err(err))
		c.JSON(http.StatusInternalServerError, BaseResponse{Error: "Failed to delete product", ErrorCode: 500})
		return
	}

	c.JSON(http.StatusOK, BaseResponse{Success: true})
}

func (h *HandlerConfig) ListProducts(c *gin.Context) {
	const op = "internal.rest.product.list"

	products, err := h.Dep.Db.ListProducts(c)
	if err != nil {
		h.Dep.Sl.Error(fmt.Sprintf("%s | Failed to list products: ", op), sl.Err(err))
		c.JSON(http.StatusInternalServerError, BaseResponse{Error: "Failed to list products", ErrorCode: 500})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}
