package apitypes

import "time"

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

func (u *User) HasPermission(requiredRole uint) bool {
	return u.AuthLevel == requiredRole
}
