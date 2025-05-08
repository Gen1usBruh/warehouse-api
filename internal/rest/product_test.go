package rest

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Gen1usBruh/warehouse-api/internal/domain/product"
	"github.com/Gen1usBruh/warehouse-api/internal/scope"
	"github.com/Gen1usBruh/warehouse-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockProductUseCase struct {
	products map[int32]product.Product
	nextID   int32
}

func (m *mockProductUseCase) Create(ctx context.Context, p product.Product) (int32, error) {
	m.nextID++
	p.ID = m.nextID
	m.products[p.ID] = p
	return p.ID, nil
}

func (m *mockProductUseCase) GetByID(ctx context.Context, id int32) (product.Product, error) {
	p, ok := m.products[id]
	if !ok {
		return product.Product{}, errors.New("not found")
	}
	return p, nil
}

func (m *mockProductUseCase) Update(ctx context.Context, p product.Product) error {
	if _, ok := m.products[p.ID]; !ok {
		return errors.New("not found")
	}
	m.products[p.ID] = p
	return nil
}

func (m *mockProductUseCase) Delete(ctx context.Context, id int32) error {
	if _, ok := m.products[id]; !ok {
		return errors.New("not found")
	}
	delete(m.products, id)
	return nil
}

func (m *mockProductUseCase) List(ctx context.Context) ([]product.Product, error) {
	var list []product.Product
	for _, p := range m.products {
		list = append(list, p)
	}
	return list, nil
}

type stubLogger struct{}

func (s *stubLogger) Error(msg string, fields ...any) {}

func setupHandlerWithMock() (*gin.Engine, *mockProductUseCase) {
	mockUC := &mockProductUseCase{
		products: make(map[int32]product.Product),
	}
	useCase := usecase.NewProductUseCase(mockUC)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	cfg := HandlerConfig{
		Dep: &scope.Dependencies{
			Product: useCase,
			Sl:      logger,
		},
	}

	return setupRouter(&cfg), mockUC
}

func setupRouter(h *HandlerConfig) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/products", h.CreateProduct)
	router.GET("/products/:id", h.GetProduct)
	router.PUT("/products/:id", h.UpdateProduct)
	router.DELETE("/products/:id", h.DeleteProduct)
	router.GET("/products", h.ListProducts)
	return router
}

func performRequest(r http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestCreateProduct(t *testing.T) {
	router, _ := setupHandlerWithMock()

	body := []byte(`{
		"name": "Test Product",
		"description": "Valid description",
		"price": 5000,
		"quantity": 10
	}`)

	resp := performRequest(router, "POST", "/products", body)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), `"id":`)
}

func TestCreateProduct_InvalidJSON(t *testing.T) {
	router, _ := setupHandlerWithMock()

	body := []byte(`{"description":"Missing name","price":10,"quantity":1}`)

	resp := performRequest(router, "POST", "/products", body)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "'Name' failed on the 'required'")
}

func TestCreateProduct_BusinessValidation(t *testing.T) {
	router, _ := setupHandlerWithMock()

	tests := []struct {
		name     string
		body     string
		expected string
	}{
		{"Price too high", `{"name":"Valid","description":"Desc","price":20000,"quantity":1}`, "price exceeds"},
		{"Reserved name", `{"name":"Sarkor","description":"Desc","price":10,"quantity":1}`, "name is reserved"},
		{"Quantity too high", `{"name":"Valid","description":"Desc","price":10,"quantity":1001}`, "quantity exceeds"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := performRequest(router, "POST", "/products", []byte(tt.body))
			assert.Equal(t, http.StatusBadRequest, resp.Code)
			assert.Contains(t, resp.Body.String(), tt.expected)
		})
	}
}

func TestGetProduct(t *testing.T) {
	router, mock := setupHandlerWithMock()

	id, _ := mock.Create(context.TODO(), product.Product{
		Name: "Iphone", Description: "Smartphone", Price: 12, Quantity: 1,
	})

	resp := performRequest(router, "GET", "/products/"+itoa(id), nil)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), `"name":"Iphone"`)
}

func TestUpdateProduct(t *testing.T) {
	router, mock := setupHandlerWithMock()

	id, _ := mock.Create(context.TODO(), product.Product{
		Name: "Old", Description: "desc", Price: 10, Quantity: 1,
	})

	body := []byte(`{
		"name": "New Name",
		"description": "Updated desc",
		"price": 22,
		"quantity": 3
	}`)
	resp := performRequest(router, "PUT", "/products/"+itoa(id), body)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestDeleteProduct(t *testing.T) {
	router, mock := setupHandlerWithMock()

	id, _ := mock.Create(context.TODO(), product.Product{
		Name: "Olcha", Description: "qizil", Price: 10, Quantity: 1,
	})

	resp := performRequest(router, "DELETE", "/products/"+itoa(id), nil)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestListProducts(t *testing.T) {
	router, mock := setupHandlerWithMock()

	mock.Create(context.TODO(), product.Product{Name: "klubnika", Description: "meva", Price: 10, Quantity: 1})
	mock.Create(context.TODO(), product.Product{Name: "pomidor", Description: "sabzavot", Price: 20, Quantity: 2})

	resp := performRequest(router, "GET", "/products", nil)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), `"name":"klubnika"`)
	assert.Contains(t, resp.Body.String(), `"name":"pomidor"`)
}

func itoa(i int32) string {
	return strconv.Itoa(int(i))
}
