package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"

	"go-chat/proto"
)

type DB struct {
	sql.DB
	userTable string
}

func OpenMySQL(dataSourceName string) (*DB, error) {
	d, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err := d.Ping(); err != nil {
		return nil, err
	}
	db := &DB{DB: *d, userTable: "users"}
	return db, nil
}

func (db *DB) Register(user *proto.User) error {
	_, err := db.Exec(`INSERT INTO users VALUES(?,?,?,?)`,
		user.ID,
		user.Name,
		user.Sex,
		user.Birth)
	return err
}

func (db *DB) IsIDExist(ID string) bool {
	row, err := db.Query("SELECT * FROM users WHERE id=?", ID)
	if err != nil {
		log.Fatal("isIDExist: %v", err)
	}
	if row.Next() {
		return true
	}
	if err := row.Err(); err != nil {
		log.Fatal("isIDExist: %v", err)
	}
	return false
}

func (db *DB) UserData(ID string) *proto.User {
	log.Print("Enter DB.UserData")
	var user proto.User
	err := db.QueryRow("SELECT * FROM users WHERE id=?", ID).Scan(
		&user.ID, &user.Name, &user.Sex, &user.Birth)
	if err != nil {
		log.Print("databse: UserData ", err)
		return nil
	}
	return &user
}

func (db *DB) PreserveMessage(data *proto.Data) error {
	_, err := db.Exec("INSERT INTO messages VALUES(?,?,?,?)",
		data.Time, data.Sender, data.Receiver, data.Content)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) RestoreMessage(ID string) []proto.Data {
	rows, err := db.Query("SELECT time, sender, receiver, content from users, messages where id=sender and receiver=?", ID)
	if err != nil {
		log.Print(err)
		return nil
	}
	defer rows.Close()
	var datas []proto.Data
	for rows.Next() {
		var data proto.Data
		if err := rows.Scan(&data.Time, &data.Sender,
			&data.Receiver, &data.Content); err != nil {
			log.Print(err)
			return nil
		}
		datas = append(datas, data)
	}
	if len(datas) == 0 {
		log.Print("none")
		return nil
	}
	if _, err := db.Exec("DELETE FROM messages WHERE receiver=?", ID); err != nil {
		log.Print(err)
	}
	return datas
}
