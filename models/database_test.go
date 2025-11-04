package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MockGormDB struct {
	mock.Mock
}

func (m *MockGormDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(dest, conds)
	return args.Get(0).(*gorm.DB)
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&Wallet{})
	require.NoError(t, err)

	return db
}

func TestSetupDataBase(t *testing.T) {
	db := setupTestDB(t)

	wallet := &Wallet{
		Person: "Джон Гарик",
		Amount: 22000,
	}

	err := db.Create(&wallet).Error

	assert.NoError(t, err)
	assert.NotZero(t, wallet.ID) // GORM должен был автоматически установить ID

	var foundWallet Wallet
	result := db.First(&foundWallet, wallet.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, wallet.Person, foundWallet.Person)
	assert.Equal(t, wallet.Amount, foundWallet.Amount)
}
