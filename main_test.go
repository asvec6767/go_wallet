package wallet

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"wallet/handlers"
	"wallet/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var testServer *handlers.Server

func initTestServer(t *testing.T) *handlers.Server {
	var err error
	db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Wallet{})
	require.NoError(t, err)

	testServer = handlers.NewServer(db)

	return testServer
}

func setupTestRouter(t *testing.T) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	testServer = initTestServer(t)

	group := router.Group("/api/v1")
	group.POST("/create", testServer.CreateWallet)
	group.POST("/wallet", testServer.WalletOperation)
	group.GET("/wallets/:wallet_uuid", testServer.WalletAmount)

	return router
}

func setupTestRouterForBench() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// db для нагрузочных тестов не связана с глобальной переменной db
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	_ = db.AutoMigrate(&models.Wallet{})

	testServer := handlers.NewServer(db)

	group := router.Group("/api/v1")
	group.POST("/create", testServer.CreateWallet)
	group.POST("/wallet", testServer.WalletOperation)
	group.GET("/wallets/:wallet_uuid", testServer.WalletAmount)

	return router
}

func TestSetupRouter_CreateWallet(t *testing.T) {
	router := setupTestRouter(t)

	wallet := models.Wallet{Person: "Джон Гарик", Amount: 22000}
	walletRequest := handlers.WalletCreateInput{Person: wallet.Person, Amount: wallet.Amount}

	jsonData, _ := json.Marshal(walletRequest)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/create", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)

	// t.Log(response)

	assert.NoError(t, err)
	assert.Equal(t, "Счет создан", response["message"])
}

func TestSetupRouter_WalletOperation(t *testing.T) {
	router := setupTestRouter(t)

	wallet := models.Wallet{Person: "Джон Гарик", Amount: 22000}
	err := db.Create(&wallet).Error
	assert.NoError(t, err)

	tests := []struct {
		name          string
		walletRequest handlers.WalletOperationInput
		shouldError   bool
	}{
		{"Проверка депозита", handlers.WalletOperationInput{WalletId: 1, OperationType: "DEPOSIT", Amount: 2000}, false},
		{"Проверка вывод", handlers.WalletOperationInput{WalletId: 1, OperationType: "WITHDRAW", Amount: 2000}, false},
		{"Несуществующая операция", handlers.WalletOperationInput{WalletId: 1, OperationType: "XXXX", Amount: 2000}, true},
	}

	for _, test := range tests {
		jsonData, _ := json.Marshal(test.walletRequest)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/v1/wallet", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if test.shouldError {
			assert.Equal(t, http.StatusBadRequest, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)

			assert.NoError(t, err)
			assert.Equal(t, "Неизвестный тип операции", response["error"])
		} else {
			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)

			assert.NoError(t, err)
			assert.Equal(t, "Операция по счету выполнена", response["message"])
		}
	}
}

func TestSetupRouter_WalletAmount(t *testing.T) {
	router := setupTestRouter(t)

	wallet := models.Wallet{Person: "Джон Гарик", Amount: 22000}
	err := db.Create(&wallet).Error
	assert.NoError(t, err)

	actualWallet, err := testServer.GetWalletById(1)
	assert.NoError(t, err)
	assert.Equal(t, wallet.Person, actualWallet.Person)
	assert.Equal(t, wallet.Amount, actualWallet.Amount)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/wallets/1", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)

	// t.Log(response)

	assert.NoError(t, err)
	assert.Equal(t, 22000, int(response["message"].(float64)))
}

func Benchmark1000RPS(b *testing.B) {
	router := setupTestRouterForBench()
	var error50xCounter int64

	b.ResetTimer()
	b.SetParallelism(1000) // 1000 параллельных горутин

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/v1/wallets/1", nil)
			router.ServeHTTP(w, req)

			if w.Code >= 500 {
				atomic.AddInt64(&error50xCounter, 1)
			}
		}
	})

	b.Logf("50X errors: %d/%d (%.2f%%)",
		atomic.LoadInt64(&error50xCounter),
		b.N,
		float64(atomic.LoadInt64(&error50xCounter))/float64(b.N)*100)
}
