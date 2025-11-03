package models

import (
	"fmt"

	"gorm.io/gorm"
)

// Модель
type Wallet struct {
	gorm.Model
	Person string `gorm:"size:255;not null" json:"person"`
	Amount int    `gorm:"not null;" json:"amount"`
}

// Внесение депозита
func (wallet *Wallet) Deposit(amount int) error {
	if amount <= 0 {
		return fmt.Errorf("Сумма депозита не может быть меньше и равной нулю")
	}

	wallet.Amount += amount

	return nil
}

// Вывод средств
func (wallet *Wallet) Withdraw(amount int) error {
	if amount <= 0 {
		return fmt.Errorf("Запрашиваемая сумма не может быть меньше или равной нулю")
	}
	if amount > wallet.Amount {
		return fmt.Errorf("Недостаточно средств на счете")
	}

	wallet.Amount -= amount

	return nil
}
