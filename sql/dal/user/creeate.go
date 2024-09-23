package user_dal

import (
	"context"
	"fmt"
	"time"

	"github.com/Arinji2/meme-backend/sql"
	"github.com/Arinji2/meme-backend/types"
)

func InitUser(ctx context.Context, user types.User) error {
	_, err := sql.ExecuteQuery(ctx, "INSERT INTO users (email, created_on) VALUES (?, ?)", user.Email, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func CreateUser(ctx context.Context, user types.User) error {
	_, err := sql.ExecuteQuery(ctx, "INSERT INTO users (username, dicebear_seed) VALUES (?, ?)", user.Username, user.DicebearSeed)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}
