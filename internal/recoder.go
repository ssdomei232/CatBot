package internal

import (
	"database/sql"
	"log"

	"git.mmeiblog.cn/mei/CatBot/pkg/napcat"
	_ "github.com/mattn/go-sqlite3"
)

func Record(groupMessage napcat.Message) {
	db, err := sql.Open("sqlite3", "history.db")
	if err != nil {
		log.Fatalf("打开数据库失败:%s", err)
	}
	defer db.Close()

	// 查询该用户是否存在于数据库中
	var id int
	sql := `SELECT id FROM user WHERE qq_number = ?`
	db.QueryRow(sql, groupMessage.Sender.UserID).Scan(&id)
	// 如果不存在，则插入该用户
	if id == 0 {
		sql = `INSERT INTO user (qq_number, nickname) VALUES (?, ?)`
		_, err = db.Exec(sql, groupMessage.Sender.UserID, groupMessage.Sender.Nickname)
		if err != nil {
			log.Fatalf("插入用户失败:%s", err)
		}
	} else {
		// 如果存在，则给用户发言数量+1
		sql = `UPDATE user SET speak_count = speak_count + 1 WHERE qq_number = ?`
		_, err = db.Exec(sql, groupMessage.Sender.UserID)
		if err != nil {
			log.Fatalf("更新用户失败:%s", err)
		}
	}

	// 记录数据
	sql = `INSERT INTO history (timestamp, group_id, group_name, sender_qq_number, sender_nickname, raw_message) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = db.Exec(sql, groupMessage.Time, groupMessage.GroupID, groupMessage.GroupName, groupMessage.Sender.UserID, groupMessage.Sender.Nickname, groupMessage.RawMessage)
	if err != nil {
		log.Fatalf("插入数据失败:%s", err)
	}
}
