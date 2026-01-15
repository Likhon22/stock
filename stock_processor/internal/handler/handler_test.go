package handler // ← SAME package as handler.go

import (
	"context"
	"sync"

	"github.com/segmentio/kafka-go"
)

// Mock the ProcessorService interface
type MockService struct {
	priceToReturn float64
	errorToReturn error
}

// Implement ALL methods from ProcessorService interface
// (We only care about GetCache, but Go requires all methods)

func (m *MockService) GetCache(ctx context.Context, symbol string) (float64, error) {
	return m.priceToReturn, m.errorToReturn
}

func (m *MockService) GetAll(ctx context.Context) (map[string]float64, error) {
	return nil, nil
}

func (m *MockService) GetHistory(ctx context.Context, symbol string, limit int) ([]float64, error) {
	return nil, nil
}

func (m *MockService) ProcessMessage(ctx context.Context, msg kafka.Message) error {
	return nil
}

func (m *MockService) StartWorkers(ctx context.Context, jobs <-chan kafka.Message, workerCount int) *sync.WaitGroup {
	return nil
}

// Now write the test
// func TestGetPrice_Success(t *testing.T) {
// 	// ARRANGE
// 	mockService := &MockService{
// 		priceToReturn: 150.45,
// 		errorToReturn: nil,
// 	}

// 	handler := &Handler{
// 		service: mockService,
// 	}

// 	// Create fake HTTP request
// 	req := httptest.NewRequest("GET", "/price/AAPL", nil)

// 	// Add URL parameter (chi router normally does this)
// 	rctx := chi.NewRouteContext()
// 	rctx.URLParams.Add("symbol", "AAPL")
// 	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

// 	// Create fake response writer
// 	w := httptest.NewRecorder()

// 	// ACT
// 	handler.GetPrice(w, req)

// 	// ASSERT
// 	// Check status code
// 	if w.Code != http.StatusOK {
// 		t.Errorf("expected status 200, got %d", w.Code)
// 	}

// 	// Parse JSON response
// 	var response map[string]interface{}
// 	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
// 		t.Fatalf("failed to decode response: %v", err)
// 	}

// 	// Check the data
// 	data, ok := response["data"].(map[string]interface{})
// 	if !ok {
// 		t.Fatal("response data is not a map")
// 	}

// 	if data["symbol"] != "AAPL" {
// 		t.Errorf("expected symbol AAPL, got %v", data["symbol"])
// 	}

// 	if data["price"] != 150.45 {
// 		t.Errorf("expected price 150.45, got %v", data["price"])
// 	}
// }
