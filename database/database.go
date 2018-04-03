package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"

	"server/chat/protocol"
)

var (
	db        *sql.DB
	userTable = "user"
)

func OpenMysql(dataSourceName string) (*sql.DB, error) {
	var err error
	db, err = sql.Open("mysql", dataSourceName)
	return db, err
}

func Register(user *protocol.User) error {
	_, err := db.Exec(`INSERT INTO users VALUES(?,?,?,?)`,
		user.ID,
		user.Name,
		user.Sex,
		user.Birth)
	return err
}

func IsIDExist(ID string) bool {
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

func UserData(ID string) *protocol.User {
	var user protocol.User
	err := db.QueryRow("SELECT * FROM users WHERE id=?", ID).Scan(
		&user.ID, &user.Name, &user.Sex, &user.Birth)
	if err != nil {
		return nil
	}
	return &user
}

func PreserveMessage(data *protocol.Data) error {
	_, err := db.Exec("INSERT INTO messages VALUES(?,?,?,?)",
		data.Time, data.ID, data.Receiver.ID, data.Content)
	if err != nil {
		return err
	}
	return nil
}

func MessagePreserved(ID string) []protocol.Data {
	rows, err := db.Query("SELECT time, sender, name, sex, birth, receiver, content from users, messages where id=sender and receiver=?", ID)
	if err != nil {
		log.Print(err)
		return nil
	}
	defer rows.Close()
	var datas []protocol.Data
	for rows.Next() {
		var data protocol.Data
		if err := rows.Scan(&data.Time,
			&data.ID, &data.Name, &data.Sex, &data.Birth,
			&data.Receiver.ID, &data.Content); err != nil {
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
