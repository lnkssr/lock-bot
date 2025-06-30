package models

type LimitRequest struct {
	User  int `json:"user_id"`
	Limit int `json:"new_limit"`
}
