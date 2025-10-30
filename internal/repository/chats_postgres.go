package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lavatee/camp_backend/internal/model"
)

type ChatsPostgres struct {
	db *sqlx.DB
}

func NewChatsPostgres(db *sqlx.DB) *ChatsPostgres {
	return &ChatsPostgres{
		db: db,
	}
}

func (r *ChatsPostgres) CreateChat(firstUserId int, secondUserId int) (int, error) {
	var id int
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("INSERT INTO %s DEFAULT VALUES RETURNING id", chatsTable)
	row := tx.QueryRow(query)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}
	query = fmt.Sprintf("INSERT INTO %s (user_id, chat_id, translating) VALUES ($1, $2, $3), ($4, $5, $6)", usersInChatTable)
	_, err = tx.Exec(query, firstUserId, id, true, secondUserId, id, true)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return 0, err
	}
	return id, nil
}

func (r *ChatsPostgres) GetUserChats(userID int, searchQuery string) ([]model.Chat, error) {
	var chats []model.Chat
	query := fmt.Sprintf(`
		SELECT 
			c.id,
			partner.id AS user_id,
			partner.name AS user_name,
			partner.photo_url AS user_photo_url,
			COALESCE(last_msg.text, '') AS last_message,
			COALESCE(last_msg.sent_at, c.created_at) AS last_message_time
		FROM %s c
		INNER JOIN %s uic1 ON c.id = uic1.chat_id AND uic1.user_id = $1
		INNER JOIN %s uic2 ON c.id = uic2.chat_id AND uic2.user_id != $1
		INNER JOIN %s partner ON uic2.user_id = partner.id
		LEFT JOIN LATERAL (
			SELECT m.text, m.sent_at
			FROM %s m
			WHERE m.chat_id = c.id
			ORDER BY m.sent_at DESC
			LIMIT 1
		) last_msg ON true
		WHERE partner.name ILIKE $2
		ORDER BY last_message_time DESC
	`, chatsTable, usersInChatTable, usersInChatTable, usersTable, messagesTable)

	err := r.db.Select(&chats, query, userID, searchQuery)
	if err != nil {
		return nil, err
	}
	return chats, nil
}

func (r *ChatsPostgres) GetOneChat(chatId int) (model.Chat, error) {
	var chat model.Chat
	query := fmt.Sprintf(`SELECT c.id, c.created_at, c.tree, u.name AS user_name, u.id AS user_id, u.photo_url AS user_photo_url, COUNT(*) AS messages_amount
	FROM %s c
	JOIN %s cu ON c.id = cu.chat_id
	JOIN %s u ON cu.user_id = u.id
	WHERE c.id = $1`, chatsTable, usersInChatTable, usersTable)
	if err := r.db.Get(&chat, query, chatId); err != nil {
		return model.Chat{}, err
	}
	return chat, nil
}

func (r *ChatsPostgres) CheckIsTreeLegit(chatId int, timeZone string) (bool, error) {
	var isTreeLegit bool
	query := fmt.Sprintf(`
        SELECT 
            last_tree_update AT TIME ZONE $2 < 
            DATE_TRUNC('day', NOW() AT TIME ZONE $2) - INTERVAL '1 day' AS is_tree_legit 
        FROM %s
        WHERE id = $1
    `, chatsTable)
	row := r.db.QueryRow(query, chatId, timeZone)
	if err := row.Scan(&isTreeLegit); err != nil {
		return false, err
	}
	return isTreeLegit, nil
}

func (r *ChatsPostgres) CreateMessage(message model.Message, timeZone string) (int, bool, error) {
	var id int
	tx, err := r.db.Begin()
	if err != nil {
		return 0, false, err
	}
	query := fmt.Sprintf("INSERT INTO %s (user_id, chat_id, text) VALUES ($1, $2, $3) RETURNING id", messagesTable)
	row := tx.QueryRow(query, message.UserId, message.ChatId, message.Text)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, false, err
	}
	isTreeUpdated, err := r.checkTree(message.UserId, message.ChatId, timeZone, tx)
	if err != nil {
		return 0, false, err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return 0, false, err
	}
	return id, isTreeUpdated, nil
}

func (r *ChatsPostgres) checkTree(userId int, chatId int, timeZone string, tx *sql.Tx) (bool, error) {
	var isUpdatedToday bool
	query := fmt.Sprintf(`
        SELECT 
            last_tree_update AT TIME ZONE $2 >= 
            DATE_TRUNC('day', NOW() AT TIME ZONE $2) AS is_updated_today
        FROM %s 
        WHERE id = $1
    `, chatsTable)
	row := tx.QueryRow(query, chatId, timeZone)
	if err := row.Scan(&isUpdatedToday); err != nil {
		return false, err
	}
	if isUpdatedToday {
		return false, nil
	}
	var isTreeMustBeUpdated bool
	query = fmt.Sprintf(`
        SELECT EXISTS(
            SELECT 1 
            FROM %s m
            JOIN %s uic ON m.chat_id = uic.chat_id
            WHERE m.chat_id = $1
            AND uic.user_id != $2
            AND m.user_id = uic.user_id
            AND m.sent_at AT TIME ZONE $3 >= DATE_TRUNC('day', NOW() AT TIME ZONE $3)
        )
    `, messagesTable, usersInChatTable)
	row = tx.QueryRow(query, chatId, userId, timeZone)
	if err := row.Scan(&isTreeMustBeUpdated); err != nil {
		return false, err
	}
	if !isTreeMustBeUpdated {
		return false, nil
	}
	query = fmt.Sprintf("UPDATE %s SET tree = tree + 1", chatsTable)
	_, err := tx.Exec(query)
	return true, err
}

func (r *ChatsPostgres) GetChatMessages(chatId int, userId int, timeZone string) ([]model.Message, error) {
	var messages []model.Message
	query := fmt.Sprintf("SELECT id, user_id, chat_id, text, sent_at AT TIME ZONE $1, is_read, user_id = $2 AS is_user_message FROM %s WHERE chat_id = $3", messagesTable)
	if err := r.db.Select(&messages, query, timeZone, userId, chatId); err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *ChatsPostgres) EditMessage(messageId int, userId int, newText string) error {
	query := fmt.Sprintf("UPDATE %s SET text = $1 WHERE id = $2 AND user_id = $3", messagesTable)
	_, err := r.db.Exec(query, newText, messageId, userId)
	return err
}

func (r *ChatsPostgres) DeleteMessage(messageId int, userId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND user_id = $2", messagesTable)
	_, err := r.db.Exec(query, messageId, userId)
	return err
}

func (r *ChatsPostgres) MakeMessageRead(messageId int) error {
	query := fmt.Sprintf("UPDATE %s SET is_read = true WHERE id = $1", messagesTable)
	_, err := r.db.Exec(query, messageId)
	return err
}
