package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

var numberId = 0

func incrementId() int {
	numberId = numberId + 1

	return numberId
}

func (s ParcelStore) Add(p *Parcel) (int, error) {
	p.Number = incrementId()
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	res, err := s.db.Exec("INSERT INTO parcel (number, client, status, address, created_at) VALUES (:number, :client, :status, :address, :created_at)",
		sql.Named("number", p.Number),
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt),
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return int(id), err
	}

	// верните идентификатор последней добавленной записи
	return int(id), nil

}

func (s ParcelStore) Get(number int) (Parcel, error) {

	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	// заполните объект Parcel данными из таблицы
	p := Parcel{}
	row := s.db.QueryRow("SELECT number, address, client, created_at, number, status FROM parcel WHERE number = :number", sql.Named("number", number))
	err := row.Scan(&p.Number, &p.Address, &p.Client, &p.CreatedAt, &p.Number, &p.Status)

	return p, err
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	var res []Parcel

	// Выполните запрос к базе данных
	row, err := s.db.Query("SELECT address, client, created_at, number, status FROM parcel WHERE client = @client", sql.Named("client", client))
	if err != nil {
		// Обработка ошибки
		return res, nil
	}
	defer row.Close()

	// Перебираем результаты запроса и заполняем массив структур
	for row.Next() {
		var p Parcel
		err := row.Scan(&p.Address, &p.Client, &p.CreatedAt, &p.Number, &p.Status)
		if err != nil {
			// Обработка ошибки
			return res, nil
		}
		res = append(res, p)
	}

	if err := row.Err(); err != nil {
		// Обработка ошибки
		return res, nil
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE Number = :number",
		sql.Named("number", number),
		sql.Named("status", status))
	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	_, err := s.db.Exec("UPDATE parcel SET address = :address",
		sql.Named("number", number),
		sql.Named("address", address))
	return err
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	//	ParcelStatusRegistered
	// удалять строку можно только если значение статуса registered
	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number AND status = 'registered'",
		sql.Named("number", number))
	return err
}
