package session_dal

import (
	"context"
	"fmt"

	"github.com/Arinji2/meme-backend/sql"
	"github.com/Arinji2/meme-backend/types"
)

func GetUserBySession(ctx context.Context, id string) (types.Session, error) {
	var session types.Session
	row, cancel := sql.ExecuteQueryRow(ctx, `
	SELECT user_id 
	FROM sessions 
	WHERE public_id = ?`, id)
	defer cancel()
	err := row.Scan(&session.UserID)
	if err != nil {
		return types.Session{}, fmt.Errorf("failed to get session by id: %w", err)
	}

	return session, nil

}
