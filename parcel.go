package main

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {

	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	// добавляем строки в таблицу parcel

	if err != nil {
		return 0, err
	}
	// возвращаем идентификатор последней добавленной записи
	id, err := res.LastInsertId()
	return int(id), err
}

func (s ParcelStore) Get(number int) (Parcel, error) {

	var (
		id        int
		client    int
		status    string
		address   string
		createdAt string
	)

	row := s.db.QueryRow("SELECT * FROM parcel WHERE number = :number", sql.Named("number", number))

	err := row.Scan(&id, &client, &status, &address, &createdAt)

	if err != nil {
		return Parcel{}, err
	}

	p := Parcel{Number: id, Client: client, Status: status, Address: address, CreatedAt: createdAt}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {

	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = :client", sql.Named("client", client))

	if err != nil {
		return []Parcel{}, err
	}
	defer rows.Close()

	var res []Parcel

	var (
		id        int
		person    int
		status    string
		address   string
		createdAt string
	)

	for rows.Next() {

		err := rows.Scan(&id, &person, &status, &address, &createdAt)

		if err != nil {
			return []Parcel{}, err
		}
		p := Parcel{Number: id, Client: person, Status: status, Address: address, CreatedAt: createdAt}
		res = append(res, p)
	}

	if err := rows.Err(); err != nil {
		return []Parcel{}, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {

	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))

	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered

	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number and status = :status",
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))

	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {

	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number and status = :status",
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return err
	}
	return nil
}
