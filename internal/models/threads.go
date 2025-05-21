package models

import (
	"database/sql"
	"time"
	"errors"
)

type Thread struct {
	Id int
	Title string
	Date time.Time
	UserId int
	UserName string
	Messages []Message
	LastMsg string
}

type ThreadModel struct{
	DB *sql.DB
}

func (tm *ThreadModel) CreateThread(threadTitle string, userId int ) (int, error) {
	stmt := `INSERT INTO Threads (title, createdDate, userId)
			VALUES (?, CURRENT_TIMESTAMP, ?);`
	result, err := tm.DB.Exec(stmt, threadTitle, userId)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (tm *ThreadModel) ShowThreadMsgs(threadId int) (*Thread, error) {
	stmt := `SELECT t.title, t.createdDate, u.userName
			FROM Threads AS t, Users AS u
			WHERE u.userId = t.userId
			AND threadId = ?;`

	row := tm.DB.QueryRow(stmt, threadId)
	var thread Thread

	messageModel := &MessageModel{DB: tm.DB}
	msgs, err := messageModel.ShowLatestMessages(threadId)
	if err != nil {
		return nil, err
	}

	err = row.Scan(&thread.Title, &thread.Date, &thread.UserName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return nil, ErrNoRecord
		}
		return nil, err
	}

	thread.Messages = msgs
	return &thread, nil
}

func (tm *ThreadModel) ShowLatestThreads() ([]Thread, error) {
	stmt := `SELECT t.threadId, t.title, t.createdDate, u.userName 
			FROM Users AS u, Threads AS t
			WHERE t.userId = u.userId 
			ORDER BY t.createdDate DESC
			LIMIT 10;`

	rows, err := tm.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var threads []Thread

	for rows.Next() {
		var thread Thread

		err = rows.Scan(&thread.Id, &thread.Title, &thread.Date, &thread.UserName)
		if err != nil {
			return nil, err
		}

		messageModel := &MessageModel{DB: tm.DB}
		lastMessage, err := messageModel.ShowLastMsg(thread.Id)
		if err != nil {
			return nil, err
		}

		if len(lastMessage.Body) > 100 {
			thread.LastMsg = lastMessage.Body[:100]
		}
		thread.LastMsg = lastMessage.Body

		threads = append(threads, thread)
		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	return threads, nil
}

func (tm *ThreadModel) IfThreadExist(title string) bool {
	var count int
	stmt := `SELECT COUNT(*)
			FROM Threads
			WHERE title = ?`
	
	row := tm.DB.QueryRow(stmt, title)
	err := row.Scan(&count)
	if err != nil || count == 0 {
		return false
	}
	return true
}

func (tm *ThreadModel) GetThreadId(title string) (int, error) {
	stmt := `SELECT threadId
			FROM Threads
			WHERE threadTitle = ?`

	var thread Thread
	row := tm.DB.QueryRow(stmt, title)
	err := row.Scan(&thread.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return 0, ErrNoRecord
		}
		return 0, err
	}
	return thread.Id, nil
} 
