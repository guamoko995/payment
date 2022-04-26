package main

import (
	"context"
	"log"
	db "payment/server/database"
	"testing"
	"time"

	proto "payment/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// UpSum (фейковый клиент) отправляет на сервис запрос на пополнение счета с переданным id
// на переданную sum с помощью gRPC.
func UpSum(client proto.PaymentClient, id, sum int64) error {
	// Установка соединения с ожиданием ответа 10 секунд.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Пополнение счета по gRPC
	_, err := client.UpSum(ctx, &proto.UpRequest{
		ID:  id,
		Sum: sum,
	})
	if err != nil {
		return err
	}
	return nil
}

// SumTransfer (фейковый клиент) отправляет на сервис запрос на перевод со счета
// senderID на счет geterID указанной sum с помощью gRPC.
func SumTransfer(client proto.PaymentClient, senderID, geterID, sum int64) error {
	// Установка соединения с ожиданием ответа 10 секунд.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Перевод средств по gRPC
	_, err := client.SumTransfer(ctx, &proto.TransferRequest{
		GeterID:  geterID,
		SenderID: senderID,
		Sum:      sum,
	})
	if err != nil {
		return err
	}
	return nil
}

// Тест работы сервиса
func TestServer(t *testing.T) {
	// Запуск сервера.
	go main()

	// Ждем, пока сервер запустится.
	time.Sleep(5 * time.Second)

	// Пподключение к базе данных сервера
	b, err := db.NewDB()
	if err != nil {
		t.Errorf("New error: %s", err)
	}

	// Добавление записей в базу данных сервера
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

	// Инициализация фейкового клиента.
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := proto.NewPaymentClient(conn)

	// Клиентское обращение к серверу с запросом на пополнение баланса
	sum := int64(1000)
	err = UpSum(client, int64(1), sum)

	if err != nil {
		t.Errorf("Set error: %s", err)
	}

	// Чтение записи базы данных сервера
	balance, err := b.GetBalance(1)
	if err != nil {
		t.Errorf("Set error: %s", err)
	}
	if balance != sum {
		t.Errorf("Expected balance: %v, added balance %v.", sum, balance)
	}

	// Клиентское обращение к серверу с запросом на перевод
	dif := int64(100)
	err = SumTransfer(client, int64(1), int64(2), dif)
	if err != nil {
		t.Errorf("Set error: %s", err)
	}

	// Чтение записей базы данных сервера
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

	// Клиентское обращение к серверу с запросом на перевод, превышающем
	// средства на счете отправителя.
	dif2 := int64(200)
	err = SumTransfer(client, int64(2), int64(1), dif2)
	expErr := "rpc error: code = Unknown desc = На счете 2 недостаточно средств"
	if err.Error() != expErr {
		t.Errorf("Expect error: %s; added error: %s", expErr, err)
	}

	// Чтение записей базы данных сервера (не должны отличаться от предыдущих)
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

	// Удаление таблицы базы данных сервера
	_, err = b.Sql.Exec("DROP TABLE accounts")
	if err != nil {
		t.Errorf("Set error: %s", err)
	}
}
