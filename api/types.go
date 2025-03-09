package apitypes

import "time"

type URLAnalytics struct {
	ShortID   string    `json:"short_id"`
	Clicks    int       `json:"clicks"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
