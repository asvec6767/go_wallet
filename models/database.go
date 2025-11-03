package models

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TODO: переделать на PostgreSQL
// Загрузка БД
func SetupDataBase() (*gorm.DB, error) {
	//Загрузка констант
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки env " + err.Error())
		// return nil, err
	}
	dbUrl := fmt.Sprint(os.Getenv("DATABASE_URL"))

	// Запуск БД
	db, err := gorm.Open(sqlite.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка при запуске БД " + err.Error())
		// return nil, err
	}

	// Включение автомиграций
	if err = db.AutoMigrate(&Wallet{}); err != nil {
		log.Fatal("Ошибка автомиграций БД " + err.Error())
		// return nil, err
	}

	return db, nil
}
