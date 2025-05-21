package models

import (
	"database/sql"
	"errors"
	"time"
)

type Message struct {
	Id int
	Body string
	Date time.Time
	UserId int
	ThreadId int
	UserName string
	ThreadTitle string
}

type MessageModel struct {
	DB *sql.DB
}

func (mm *MessageModel) CreateMsg(msgBody string, userId int, threadId int) (int, error){
	stmt := `INSERT INTO Messages (msgBody, createdDate, userId, threadId)
			VALUES (?, CURRENT_TIMESTAMP, ?, ?);`

	result, err := mm.DB.Exec(stmt, msgBody, userId, threadId)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (mm *MessageModel) ShowMsg(threadId int, msgId int) (*Message, error) {
	stmt := `SELECT m.msgBody, m.createdDate, u.userName, t.title
			FROM Messages AS m, Users AS u, Threads AS t
			WHERE u.userId = m.userId 
			AND m.threadId = t.threadId
			AND m.threadId = ?
			AND m.msgId = ?;`
	
	row := mm.DB.QueryRow(stmt, threadId, msgId)
	var msg Message

	err := row.Scan(&msg.Body, &msg.Date, &msg.UserName, &msg.ThreadTitle)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &msg, nil
}

func (mm *MessageModel) ShowLastMsg(threadId int) (*Message, error){
	stmt := `SELECT m.msgBody, m.createdDate, u.userName, t.title
			FROM Users AS u, Threads AS t, Messages AS m
			WHERE m.userId = u.userId 
			AND t.threadId = m.threadId
			AND t.threadId = ?
			ORDER BY m.createdDate DESC
			LIMIT 1;`

	row := mm.DB.QueryRow(stmt, threadId)
	var msg Message

	err := row.Scan(&msg.Body, &msg.Date, &msg.UserName, &msg.ThreadTitle)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &msg, nil
}

func (mm * MessageModel) ShowLatestMessages(threadId int) ([]Message, error) {
	stmt := `SELECT m.msgBody, m.createdDate, u.userName, t.title
			FROM Users AS u, Threads AS t, Messages AS m
			WHERE m.userId = u.userId 
			AND t.threadId = m.threadId
			AND t.threadId = ?
			ORDER BY m.createdDate DESC
			LIMIT 10;`

	rows, err := mm.DB.Query(stmt, threadId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []Message

	for rows.Next() {
		var msg Message
		err = rows.Scan(&msg.Body, &msg.Date, &msg.UserName, &msg.ThreadTitle)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	return messages, nil
}
