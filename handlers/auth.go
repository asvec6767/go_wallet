package handlers

import (
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Форма отправки операции
type WalletOperationInput struct {
	WalletId      int    `json:"walletid" binding:"required"`
	OperationType string `json:"operationtype" binding:"required"`
	Amount        int    `json:"amount" binding:"required"`
}

// Форма отправки операции
type WalletCreateInput struct {
	Person string `json:"person" binding:"required"`
	Amount int    `json:"amount" binding:"required"`
}

// Сервер
type Server struct {
	db *gorm.DB
}

func NewServer(db *gorm.DB) *Server {
	return &Server{db: db}
}

func (server *Server) CreateWallet(ctx *gin.Context) {
	var input WalletCreateInput

	// Бинд формы с моделью при ее получении с фронта
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка обработки формы " + err.Error()})
		return
	}

	// Запись ввода в модель пользователя
	user := models.Wallet{Person: input.Person, Amount: input.Amount}

	// Запись модели в БД
	if err := server.db.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка записи в БД " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Счет создан"})
}

func (server *Server) WalletOperation(ctx *gin.Context) {
	var input WalletOperationInput

	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := server.GetWalletById(uint(input.WalletId))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch input.OperationType {
	case "DEPOSIT":
		if err = wallet.Deposit(input.Amount); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	case "WITHDRAW":
		if err = wallet.Withdraw(input.Amount); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неизвестный тип операции"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Операция по счету выполнена"})
}

func (server *Server) GetWalletById(id uint) (models.Wallet, error) {
	wallet := models.Wallet{}

	if err := server.db.Model(models.Wallet{}).Where("id=?", id).Take(&wallet).Error; err != nil {
		return wallet, err
	}

	return wallet, nil
}
