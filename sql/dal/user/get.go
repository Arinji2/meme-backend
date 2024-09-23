package user_dal

import (
	"context"
	"fmt"

	"github.com/Arinji2/meme-backend/sql"
	"github.com/Arinji2/meme-backend/types"
)

func GetUserByEmail(ctx context.Context, email string) (types.User, error) {
	var user types.User
	userID := ctx.Value(types.UserIDKey)
	fmt.Println(string(userID.([]uint8)))
	row, cancel := sql.ExecuteQueryRow(ctx, "SELECT id, username, email, dicebear_seed, created_on FROM users WHERE email = ?", email)
	defer cancel()
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.DicebearSeed, &user.CreatedOn)
	if err != nil {
		return types.User{}, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}
