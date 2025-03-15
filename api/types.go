package apitypes

import (
	"time"
)

type URL struct {
	ID          string    `json:"id"`
	ShortID     string    `json:"short_id"`
	OriginalURL string    `json:"original_url"`
	ExpiresAt   time.Time `json:"expires_at"`
	UserID      uint      `json:"user_id"`
}

type URLAnalytics struct {
	ShortID   string    `json:"short_id"`
	Clicks    int       `json:"clicks"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	AuthLevel uint      `json:"auth_level"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) FreeUser() bool {
	return u.AuthLevel == 1
}

func (u *User) PremiumUser() bool {
	return u.AuthLevel == 2
}

func (u *User) Admin() bool {
	return u.AuthLevel == 3
}
