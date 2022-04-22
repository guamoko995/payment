package database

import (
	"testing"
)

// Тест работы базы данных
func TestDatabase(t *testing.T) {
	// Создание временной базы данных (с пустым именем файла)
	b, err := NewDB("")
	if err != nil {
		t.Errorf("New error: %s", err)
	}

	// Добавление записей
	// №1
	err = b.AddAccount()
	if err != nil {
		t.Errorf("Add error: %s", err)
	}
	// №2
	err = b.AddAccount()
	if err != nil {
		t.Errorf("Add error: %s", err)
	}

	// Пополнение баланса
	sum := int64(1000)
	err = b.UpBalance(1, sum)
	if err != nil {
		t.Errorf("Set error: %s", err)
	}

	// Чтение записи
	balance, err := b.GetBalance(1)
	if err != nil {
		t.Errorf("Set error: %s", err)
	}
	if balance != sum {
		t.Errorf("Expected balance: %v, added balance %v.", balance, sum)
	}

	// Списание с баланса
	decr := int64(500)
	err = b.DownBalance(1, decr)
	if err != nil {
		t.Errorf("Set error: %s", err)
	}

	// Чтение записи
	balance, err = b.GetBalance(1)
	if err != nil {
		t.Errorf("Set error: %s", err)
	}
	sum -= decr
	if balance != sum {
		t.Errorf("Expected balance: %v, added balance %v.", balance, sum)
	}

	// Перевод
	dif := int64(100)
	err = b.AmountTransfer(1, 2, dif)
	if err != nil {
		t.Errorf("Set error: %s", err)
	}

	// Чтение записей
	//#1
	balance, err = b.GetBalance(1)
	if err != nil {
		t.Errorf("Set error: %s", err)
	}
	if balance != sum-dif {
		t.Errorf("Expected balance: %v, added balance %v.", sum-dif, balance)
	}
	//#2
	balance, err = b.GetBalance(2)
	if err != nil {
		t.Errorf("Set error: %s", err)
	}
	if balance != dif {
		t.Errorf("Expected balance: %v, added balance %v.", dif, balance)
	}
}
