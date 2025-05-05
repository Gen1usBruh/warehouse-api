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

// CreateProduct godoc
// @Summary Create a new product
// @Description Add a new product to the warehouse
// @Tags products
// @Accept json
// @Produce json
// @Param product body ProductRequest true "Product info"
// @Success 200 {object} map[string]int "Returns ID of created product"
// @Failure 400 {object} BaseResponse "Invalid input"
// @Failure 500 {object} BaseResponse "Server error"
// @Router /products [post]
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

// GetProduct godoc
// @Summary Get product by ID
// @Description Retrieve a single product from the warehouse
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]db.GetProductByIDRow "Product data"
// @Failure 400 {object} BaseResponse "Invalid ID"
// @Failure 404 {object} BaseResponse "Product not found"
// @Router /products/{id} [get]
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

// UpdateProduct godoc
// @Summary Update product by ID
// @Description Update product information in the warehouse
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body ProductRequest true "Updated product info"
// @Success 200 {object} BaseResponse "Product updated"
// @Failure 400 {object} BaseResponse "Invalid input"
// @Failure 500 {object} BaseResponse "Update failed"
// @Router /products/{id} [put]
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

// DeleteProduct godoc
// @Summary Delete product by ID
// @Description Remove a product from the warehouse
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} BaseResponse "Product deleted"
// @Failure 400 {object} BaseResponse "Invalid ID"
// @Failure 500 {object} BaseResponse "Delete failed"
// @Router /products/{id} [delete]
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

// ListProducts godoc
// @Summary List all products
// @Description Get a list of all products in the warehouse
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} map[string][]db.ListProductsRow "List of products"
// @Failure 500 {object} BaseResponse "List retrieval failed"
// @Router /products [get]
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
