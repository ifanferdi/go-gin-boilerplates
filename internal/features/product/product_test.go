package product

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gin-boilerplate/internal/config"
	"go-gin-boilerplate/internal/features/category"
	services "go-gin-boilerplate/internal/features/product/_services"
	grpcinfra "go-gin-boilerplate/internal/infrastructure/grpc"
	"go-gin-boilerplate/internal/middlewares"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/structpb"
)

type mockCategoryStore struct {
	items  []category.Category
	nextID uint
}

func newMockCategoryStore() *mockCategoryStore {
	return &mockCategoryStore{items: []category.Category{}, nextID: 1}
}

func (m *mockCategoryStore) FindAll(_ category.FindAllCategoryDTO) ([]category.Category, error) {
	return m.items, nil
}
func (m *mockCategoryStore) Count(_ category.FindAllCategoryDTO) (int64, error) {
	return int64(len(m.items)), nil
}
func (m *mockCategoryStore) FindByID(id uint) (*category.Category, error) {
	for i := range m.items {
		if m.items[i].ID == id {
			return &m.items[i], nil
		}
	}
	return nil, nil
}
func (m *mockCategoryStore) FindByIDs(ids []uint) ([]category.Category, error) {
	var out []category.Category
	for _, id := range ids {
		for _, c := range m.items {
			if c.ID == id {
				out = append(out, c)
			}
		}
	}
	return out, nil
}
func (m *mockCategoryStore) CountByName(name string, excludeID uint) (int64, error) {
	var c int64
	for _, it := range m.items {
		if it.Name == name && it.ID != excludeID {
			c++
		}
	}
	return c, nil
}
func (m *mockCategoryStore) CreateRecord(name string, description *string) (*category.Category, error) {
	cat := category.Category{ID: m.nextID, Name: name, Description: description}
	m.nextID++
	m.items = append(m.items, cat)
	return &cat, nil
}
func (m *mockCategoryStore) Create(cat *category.Category) error {
	cat.ID = m.nextID
	m.nextID++
	m.items = append(m.items, *cat)
	return nil
}
func (m *mockCategoryStore) Update(cat *category.Category) error {
	for i := range m.items {
		if m.items[i].ID == cat.ID {
			m.items[i] = *cat
		}
	}
	return nil
}
func (m *mockCategoryStore) Delete(id uint) error {
	for i := range m.items {
		if m.items[i].ID == id {
			m.items = append(m.items[:i], m.items[i+1:]...)
			return nil
		}
	}
	return nil
}

type mockProductRepo struct {
	catStore *mockCategoryStore
	items    []Product
	pivots   []ProductHasCategory
	nextID   uint
}

func newMockProductRepo(catStore *mockCategoryStore) *mockProductRepo {
	return &mockProductRepo{catStore: catStore, items: []Product{}, pivots: []ProductHasCategory{}, nextID: 1}
}

func (m *mockProductRepo) FindAll(_ FindAllProductDTO) ([]Product, error) {
	return m.items, nil
}
func (m *mockProductRepo) Count(_ FindAllProductDTO) (int64, error) {
	return int64(len(m.items)), nil
}
func (m *mockProductRepo) FindByID(id uint) (*Product, error) {
	for i := range m.items {
		if m.items[i].ID == id {
			return &m.items[i], nil
		}
	}
	return nil, nil
}
func (m *mockProductRepo) CountBySKU(sku string, excludeID uint) (int64, error) {
	var c int64
	for _, it := range m.items {
		if it.SKU != nil && *it.SKU == sku && it.ID != excludeID {
			c++
		}
	}
	return c, nil
}
func (m *mockProductRepo) CountByBarcode(barcode string, excludeID uint) (int64, error) {
	var c int64
	for _, it := range m.items {
		if it.Barcode != nil && *it.Barcode == barcode && it.ID != excludeID {
			c++
		}
	}
	return c, nil
}
func (m *mockProductRepo) CountExistingCategories(ids []uint) (int64, error) {
	cats, _ := m.catStore.FindByIDs(ids)
	return int64(len(cats)), nil
}
func (m *mockProductRepo) CreateWithCategories(input services.CreateProductInput) (*Product, error) {
	cats, _ := m.catStore.FindByIDs(input.CategoryIDs)
	p := Product{
		ID:         m.nextID,
		Name:       input.Name,
		Price:      input.Price,
		SKU:        input.SKU,
		Barcode:    input.Barcode,
		Stock:      input.Stock,
		IsActive:   true,
		Categories: cats,
	}
	m.nextID++
	for _, c := range cats {
		m.pivots = append(m.pivots, ProductHasCategory{ProductID: p.ID, CategoryID: c.ID})
	}
	m.items = append(m.items, p)
	return &p, nil
}
func (m *mockProductRepo) UpdateWithCategories(item *Product, categoryIDs *[]uint) (*Product, error) {
	if categoryIDs != nil {
		cats, _ := m.catStore.FindByIDs(*categoryIDs)
		item.Categories = cats
	}
	for i := range m.items {
		if m.items[i].ID == item.ID {
			m.items[i] = *item
		}
	}
	return item, nil
}
func (m *mockProductRepo) Delete(id uint) error {
	for i := range m.items {
		if m.items[i].ID == id {
			m.items = append(m.items[:i], m.items[i+1:]...)
			return nil
		}
	}
	return nil
}

func TestCategoryAndProductManyToManyFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	catRepo := newMockCategoryStore()
	catSvc := category.NewCategoryService(catRepo)

	prodRepo := newMockProductRepo(catRepo)
	prodSvc := NewProductService(prodRepo)

	r := gin.New()
	r.Use(middlewares.ErrorHandler(&config.Config{App: config.AppConfig{Env: "test"}}))
	v1 := r.Group("/api/v1")
	category.RegisterHTTPRoutes(v1, catSvc)
	RegisterHTTPRoutes(v1, prodSvc)

	// 1. Create Category via HTTP (delegates to category/_services/category_service_create.go)
	catBody, _ := json.Marshal(map[string]string{"name": "Coffee", "description": "Espresso drinks"})
	reqCat := httptest.NewRequest(http.MethodPost, "/api/v1/categories", bytes.NewReader(catBody))
	reqCat.Header.Set("Content-Type", "application/json")
	wCat := httptest.NewRecorder()
	r.ServeHTTP(wCat, reqCat)
	if wCat.Code != http.StatusCreated {
		t.Fatalf("expected 201 creating category, got %d: %s", wCat.Code, wCat.Body.String())
	}

	// 2. Create Product linked to Category ID 1 via pivot product_has_categories
	// (delegates to product/_services/product_service_create.go)
	prodBody, _ := json.Marshal(map[string]interface{}{
		"name":        "Iced Americano",
		"price":       25000,
		"sku":         "SKU-COFFEE-001",
		"stock":       50,
		"categoryIds": []uint{1},
	})
	reqProd := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(prodBody))
	reqProd.Header.Set("Content-Type", "application/json")
	wProd := httptest.NewRecorder()
	r.ServeHTTP(wProd, reqProd)
	if wProd.Code != http.StatusCreated {
		t.Fatalf("expected 201 creating product, got %d: %s", wProd.Code, wProd.Body.String())
	}
	if len(prodRepo.pivots) != 1 || prodRepo.pivots[0].CategoryID != 1 || prodRepo.pivots[0].ProductID != 1 {
		t.Fatalf("expected pivot product_has_categories entry {ProductID:1, CategoryID:1}, got %+v", prodRepo.pivots)
	}

	// 3. Invalid categoryId returns 400 Bad Request
	badProdBody, _ := json.Marshal(map[string]interface{}{
		"name":        "Latte",
		"price":       30000,
		"categoryIds": []uint{999},
	})
	reqBad := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(badProdBody))
	reqBad.Header.Set("Content-Type", "application/json")
	wBad := httptest.NewRecorder()
	r.ServeHTTP(wBad, reqBad)
	if wBad.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for non-existent categoryId, got %d", wBad.Code)
	}

	// 4. Verify gRPC Server & Client Connection for both CategoryService and ProductService
	grpcSrv := grpcinfra.NewServer()
	category.RegisterGRPCRoutes(grpcSrv.Engine, catSvc)
	RegisterGRPCRoutes(grpcSrv.Engine, prodSvc)
	defer grpcSrv.Stop()

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	go func() { _ = grpcSrv.Engine.Serve(lis) }()

	conn, err := grpcinfra.NewClientConnection(lis.Addr().String())
	if err != nil {
		t.Fatalf("failed to connect gRPC client: %v", err)
	}
	defer conn.Close()

	grpcReq, _ := structpb.NewStruct(map[string]interface{}{"id": 1})
	grpcRes := new(structpb.Struct)
	if err := conn.Invoke(context.Background(), "/product.v1.ProductService/FindByID", grpcReq, grpcRes); err != nil {
		t.Fatalf("unexpected gRPC Invoke error: %v", err)
	}
	if got := grpcRes.GetFields()["name"].GetStringValue(); got != "Iced Americano" {
		t.Fatalf("expected gRPC product name 'Iced Americano', got %q", got)
	}
}
