package types

import "time"

type GoogleUserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Provider struct {
	ID           []uint8
	UserID       []uint8
	ProviderID   int
	RefreshToken string
	AccessToken  string
	ExpiresOn    time.Time
}
