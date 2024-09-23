package types

import "time"

type Session struct {
	ID        []uint8
	PublicID  []uint8
	UserID    []uint8
	ExpiresOn time.Time
}
