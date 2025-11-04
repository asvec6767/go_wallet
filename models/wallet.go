package models

import (
	"errors"

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
		return errors.New("сумма депозита не может быть меньше и равной нулю")
	}

	wallet.Amount += amount

	return nil
}

// Вывод средств
func (wallet *Wallet) Withdraw(amount int) error {
	if amount <= 0 {
		return errors.New("запрашиваемая сумма не может быть меньше или равной нулю")
	}
	if amount > wallet.Amount {
		return errors.New("недостаточно средств на счете")
	}

	wallet.Amount -= amount

	return nil
}
