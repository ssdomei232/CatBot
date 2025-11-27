package internal

import (
	"database/sql"
	"log"

	"git.mmeiblog.cn/mei/CatBot/handler"
	napcat "github.com/ssdomei232/go-napcat-ws"
	_ "modernc.org/sqlite"
)

var messageQueue = make(chan napcat.Message, 1000)

func init() {
	go processMessages()
}

func processMessages() {
	for message := range messageQueue {
		recordMessage(message)
	}
}

func Record(groupMessage napcat.Message) {
	select {
	case messageQueue <- groupMessage:
	default:
		log.Println("消息队列已满，丢弃消息")
	}
}

func recordMessage(groupMessage napcat.Message) {
	db, err := handler.GetDB()
	if err != nil {
		log.Fatalf("打开数据库失败:%s", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	defer db.Close()

	// 查询该用户是否存在于数据库中
	var id int
	query := `SELECT id FROM user WHERE qq_number = ?`
	err = db.QueryRow(query, groupMessage.Sender.UserID).Scan(&id)

	// 如果不存在，则插入该用户
	if err == sql.ErrNoRows || id == 0 {
		query = `INSERT INTO user (qq_number, nickname, speak_count) VALUES (?, ?, 1)`
		_, err = db.Exec(query, groupMessage.Sender.UserID, groupMessage.Sender.Nickname)
		if err != nil {
			log.Fatalf("插入用户失败:%s", err)
		}
	} else if err != nil {
		log.Fatalf("查询用户失败:%s", err)
	} else {
		// 如果存在，则给用户发言数量+1
		query = `UPDATE user SET speak_count = speak_count + 1 WHERE qq_number = ?`
		_, err = db.Exec(query, groupMessage.Sender.UserID)
		if err != nil {
			log.Fatalf("更新用户失败:%s", err)
		}
	}

	// 记录数据
	query = `INSERT INTO history (timestamp, group_id, group_name, sender_qq_number, sender_nickname, raw_message) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = db.Exec(query, groupMessage.Time, groupMessage.GroupID, groupMessage.GroupName, groupMessage.Sender.UserID, groupMessage.Sender.Nickname, groupMessage.RawMessage)
	if err != nil {
		log.Fatalf("插入数据失败:%s", err)
	}
}
