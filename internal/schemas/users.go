package schemas

import "time"

type User struct {
	ID             int64     `json:"id"`
	Email          string    `json:"email"`
	Plan           string    `json:"plan"`
	PlanBoughtDate time.Time `json:"plan_bought_date"`
	MessagesLeft   int       `json:"messages_left"`
	BytesDataLeft  int       `json:"bytes_data_left"`
	BotsLeft       int       `json:"bots_left"`
}

type SaveUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Plan     string `json:"plan"`
}

type SaveUserResponse struct {
	UserID int64 `json:"id"`
}

type UpdatePlanRequest struct {
	Plan string `json:"plan"`
}

type UpdatePlanResponse struct {
	Success bool `json:"success"`
}
