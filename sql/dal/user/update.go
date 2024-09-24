package user_dal

import (
	"context"
	"fmt"

	sql_db "github.com/Arinji2/meme-backend/sql"
	"github.com/Arinji2/meme-backend/types"
)

func UpdateUserName(ctx context.Context, user types.User) error {
	_, err := sql_db.ExecuteQuery(ctx, "UPDATE users SET username = ? WHERE id = ?", user.Username, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}
