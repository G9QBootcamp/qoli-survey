package dto

import "time"

type GetNotificationsRequest struct {
	Page int  `query:"page"`
	Seen bool `query:"seen"`
}

type NotificationResponse struct {
	ID        uint      `json:"id"`
	Message   string    `json:"message"`
	Seen      bool      `json:"seen"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
