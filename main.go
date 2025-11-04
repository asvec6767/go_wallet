package main

import (
	"log"
	"main/handlers"
	"main/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// Инициализация БД в main.go
func DBInit() *gorm.DB {
	db, err := models.SetupDataBase()
	if err != nil {
		log.Println("Проблема при загрузке БД")
	}

	return db
}

func SetupRouter() *gin.Engine {
	//Создание роутера
	router := gin.Default()

	// Инициализация БД
	db := DBInit()
	server := handlers.NewServer(db)

	// Подключение шаблонов
	router.LoadHTMLGlob("templates/*")

	// Маршруты
	router.GET("/edit", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "edit.html", gin.H{})
	})

	group := router.Group("/api/v1")
	group.POST("/create", server.CreateWallet)
	group.POST("/wallet", server.WalletOperation)
	group.POST("/wallets/:wallet_uuid", server.WalletAmount)

	return router
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Ошибка при заугрзке файла .env")
	}
	port := os.Getenv("PORT")

	router := SetupRouter()

	//Запуск сервера
	// log.Println("Сервер запущен на http://localhost:8080/")
	log.Fatal(router.Run(":" + port))
}
