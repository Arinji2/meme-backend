package user_dal

import (
	"context"
	"fmt"
	"time"

	custom_log "github.com/Arinji2/meme-backend/logger"
	sql_db "github.com/Arinji2/meme-backend/sql"
	"github.com/Arinji2/meme-backend/types"
)

func InitUser(ctx context.Context, user types.User) (types.User, error) {
	start := time.Now()
	defer func() {
		custom_log.Logger.Infof("InitUser took %v to execute", time.Since(start))
	}()

	var query string
	var args []interface{}

	if user.Username == "" {
		query = "INSERT INTO users (email, created_on) VALUES (?, ?) RETURNING id, username, email, dicebear_seed, created_on"
		args = []interface{}{user.Email, time.Now().Format("2006-01-02 15:04:05")}
	} else {
		query = "INSERT INTO users (email, username, created_on) VALUES (?, ?, ?) RETURNING id, username, email, dicebear_seed, created_on"
		args = []interface{}{user.Email, user.Username, time.Now().Format("2006-01-02 15:04:05")}
	}

	rows := sql_db.ExecuteQueryRow(ctx, query, args...)

	var createdUser types.User
	err := rows.Scan(
		&createdUser.ID,
		&createdUser.Username,
		&createdUser.Email,
		&createdUser.DicebearSeed,
		&createdUser.CreatedOn,
	)

	if err != nil {
		return types.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, nil
}
