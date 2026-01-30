package service

import "time"

type Users struct {
	Id         int       `json:"id"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	Role       string    `json:"role"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
