package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithdraw(t *testing.T) {
	tests := []struct {
		name        string
		object      Wallet
		input       int
		shouldError bool
	}{
		{"снятие части денег", Wallet{Person: "Бейби Мело", Amount: 35000}, 13000, false},
		{"снятие нуля", Wallet{Person: "Джон Гарик", Amount: 15000}, 0, true},
		{"снятие всех денег", Wallet{Person: "Воскресенский", Amount: 7000}, 7000, false},
		{"снятие больше максимума", Wallet{Person: "Молодой Калуга", Amount: 7000}, 7500, true},
		{"снятие отрицательной суммы", Wallet{Person: "Поляна", Amount: 9000}, -2000, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.object.Withdraw(test.input)

			if test.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeposit(t *testing.T) {
	tests := []struct {
		name        string
		object      Wallet
		input       int
		shouldError bool
	}{
		{"Добавление суммы", Wallet{Person: "Кореш", Amount: 24000}, 6000, false},
		{"Добавление нулевой суммы", Wallet{Person: "Эксайл", Amount: 15000}, 0, true},
		{"Добавление отрицательной суммы", Wallet{Person: "Парадеич", Amount: 18000}, -3000, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.object.Deposit(test.input)

			if err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
