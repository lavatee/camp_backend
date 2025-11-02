package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lavatee/camp_backend/internal/model"
)

type RoomsPostgres struct {
	db *sqlx.DB
}

func NewRoomsPostgres(db *sqlx.DB) *RoomsPostgres {
	return &RoomsPostgres{
		db: db,
	}
}

func (r *RoomsPostgres) JoinRoom(userId int) (model.Room, error) {
    var targetRoomId int
    findQuery := fmt.Sprintf(`
        SELECT room_id 
        FROM %s 
        GROUP BY room_id 
        HAVING COUNT(*) = 1 
        LIMIT 1
        FOR UPDATE SKIP LOCKED
    `, usersInRoomTable)
    err := r.db.Get(&targetRoomId, findQuery)
    if err == nil {
        insertQuery := fmt.Sprintf(`
            INSERT INTO %s (user_id, room_id) 
            VALUES ($1, $2)
            RETURNING room_id
        `, usersInRoomTable)
        
        err = r.db.Get(&targetRoomId, insertQuery, userId, targetRoomId)
        if err != nil {
            return model.Room{}, err
        }
        
        return r.getRoomWithUsers(targetRoomId, userId)
    }
    return r.createNewRoom(userId)
}

func (r *RoomsPostgres) createNewRoom(userId int) (model.Room, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return model.Room{}, err
	}
	defer tx.Rollback()
	var room model.Room
	query := fmt.Sprintf("INSERT INTO %s DEFAULT VALUES RETURNING id", roomsTable)
	err = tx.Get(&room.Id, query)
	if err != nil {
		return model.Room{}, err
	}
	query = fmt.Sprintf("INSERT INTO %s (user_id, room_id) VALUES ($1, $2)", usersInRoomTable)
	_, err = tx.Exec(query, userId, room.Id)
	if err != nil {
		return model.Room{}, err
	}
	if err = tx.Commit(); err != nil {
		return model.Room{}, err
	}
	return room, nil
}

func (r *RoomsPostgres) getRoomWithUsers(roomId int, currentUserId int) (model.Room, error) {
	var room model.Room
	query := fmt.Sprintf(`
        SELECT
            r.id,
            u.id AS user_id,
            u.name AS user_name,
            u.photo_url AS user_photo_url
        FROM %s r
        JOIN %s ru ON r.id = ru.room_id
        JOIN %s u ON u.id = ru.user_id
        WHERE r.id = $1 
        AND u.id != $2
        LIMIT 1
    `, roomsTable, usersInRoomTable, usersTable)

	err := r.db.Get(&room, query, roomId, currentUserId)
	return room, err
}

func (r *RoomsPostgres) LeaveRoom(userId int, roomId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND room_id = $2", usersInRoomTable)
	_, err := r.db.Exec(query, userId, roomId)
	return err
}

func (r *RoomsPostgres) GetRoomUser(userId int, roomId int) (model.Room, error) {
	var roomUser model.Room
	query := fmt.Sprintf("SELECT ru.translating, u.language, u.name AS user_name, u.id AS user_id, r.id, u.photo_url AS user_photo_url FROM %s r JOIN %s ru ON r.id = ru.room_id JOIN %s u ON u.id = ru.user_id WHERE r.id = $1 AND u.id != $2", roomsTable, usersInRoomTable, usersTable)
	if err := r.db.Get(&roomUser, query, roomId, userId); err != nil {
		return model.Room{}, err
	}
	return roomUser, nil
}
