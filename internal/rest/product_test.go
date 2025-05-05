package rest

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	db "github.com/Gen1usBruh/warehouse-api/internal/storage/postgres/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockDB struct {
	products map[int32]db.Product
	nextID   int32
}

func (m *mockDB) CreateProduct(ctx context.Context, arg db.CreateProductParams) (int32, error) {
	m.nextID++
	p := db.Product{
		ID:          m.nextID,
		Name:        arg.Name,
		Description: arg.Description,
		Price:       arg.Price,
		Quantity:    arg.Quantity,
	}
	m.products[p.ID] = p
	return p.ID, nil
}

func (m *mockDB) GetProductByID(ctx context.Context, id int32) (db.Product, error) {
	p, ok := m.products[id]
	if !ok {
		return db.Product{}, errors.New("not found")
	}
	return p, nil
}

func (m *mockDB) UpdateProduct(ctx context.Context, arg db.UpdateProductParams) error {
	if _, ok := m.products[arg.ID]; !ok {
		return errors.New("not found")
	}
	m.products[arg.ID] = db.Product{
		ID:          arg.ID,
		Name:        arg.Name,
		Description: arg.Description,
		Price:       arg.Price,
		Quantity:    arg.Quantity,
	}
	return nil
}

func (m *mockDB) DeleteProduct(ctx context.Context, id int32) error {
	if _, ok := m.products[id]; !ok {
		return errors.New("not found")
	}
	delete(m.products, id)
	return nil
}

func (m *mockDB) ListProducts(ctx context.Context) ([]db.Product, error) {
	var list []db.Product
	for _, p := range m.products {
		list = append(list, p)
	}
	return list, nil
}

// Stub logger
type stubLogger struct{}

func (s *stubLogger) Error(msg string, fields ...any) {}

func setupHandler() *HandlerConfig {
	return &HandlerConfig{
		Dep: Deps{
			Db: &mockDB{products: map[int32]db.Product{}, nextID: 0},
			Sl: &stubLogger{},
		},
	}
}

func performRequest(r http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestCreateProduct(t *testing.T) {
	h := setupHandler()
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/products", h.CreateProduct)

	body := []byte(`{
		"name": "Test Product",
		"description": "Some description",
		"price": 19,
		"quantity": 5
	}`)
	w := performRequest(router, "POST", "/products", body)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"id":`)
}

func TestGetProduct(t *testing.T) {
	h := setupHandler()
	db := h.Dep.Db.(*mockDB)
	id, _ := db.CreateProduct(context.TODO(), db.CreateProductParams{
		Name: "Test", Description: "desc", Price: 12, Quantity: 1,
	})

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/products/:id", h.GetProduct)

	w := performRequest(router, "GET", "/products/"+strconv.Itoa(int(id)), nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"name":"Test"`)
}

func TestUpdateProduct(t *testing.T) {
	h := setupHandler()
	db := h.Dep.Db.(*mockDB)
	id, _ := db.CreateProduct(context.TODO(), db.CreateProductParams{
		Name: "Old", Description: "desc", Price: 10, Quantity: 1,
	})

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.PUT("/products/:id", h.UpdateProduct)

	body := []byte(`{
		"name": "New Name",
		"description": "Updated desc",
		"price": 22,
		"quantity": 3
	}`)
	w := performRequest(router, "PUT", "/products/"+strconv.Itoa(int(id)), body)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteProduct(t *testing.T) {
	h := setupHandler()
	db := h.Dep.Db.(*mockDB)
	id, _ := db.CreateProduct(context.TODO(), db.CreateProductParams{
		Name: "Del", Description: "desc", Price: 10, Quantity: 1,
	})

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.DELETE("/products/:id", h.DeleteProduct)

	w := performRequest(router, "DELETE", "/products/"+strconv.Itoa(int(id)), nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListProducts(t *testing.T) {
	h := setupHandler()
	db := h.Dep.Db.(*mockDB)
	_, _ = db.CreateProduct(context.TODO(), db.CreateProductParams{
		Name: "P1", Description: "desc", Price: 10, Quantity: 1,
	})
	_, _ = db.CreateProduct(context.TODO(), db.CreateProductParams{
		Name: "P2", Description: "desc", Price: 20, Quantity: 2,
	})

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/products", h.ListProducts)

	w := performRequest(router, "GET", "/products", nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"name":"P1"`)
	assert.Contains(t, w.Body.String(), `"name":"P2"`)
}
