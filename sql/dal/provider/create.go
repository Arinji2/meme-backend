package provider_dal

import (
	"context"
	"fmt"

	sql_db "github.com/Arinji2/meme-backend/sql"
	"github.com/Arinji2/meme-backend/types"
)

func CreateProvider(ctx context.Context, user types.User, provider types.Provider) error {
	_, err := sql_db.ExecuteQuery(ctx, "INSERT INTO oauth_providers (user_id, provider_id, refresh_token, access_token, expires_on) VALUES (?, ?, ?, ?, ?)", user.ID, provider.ProviderID, provider.RefreshToken, provider.AccessToken, provider.ExpiresOn)
	if err != nil {
		return fmt.Errorf("failed to create provider for user %s: %w", user.Email, err)
	}
	return nil
}
