package user

import (
	resp "github.com/AlexandrLitkevich/qwery/internal/lib/api/response"
)

type User struct {
	ID       string `json:"ID" validate:"required,ID"`
	Name     string `json:"name"`
	Age      int    `json:"age" `
	Position string `json:"position"`
}

type Request struct {
	Name     string `json:"name"`
	Age      int    `json:"age" `
	Position string `json:"position"`
}

type Response struct {
	resp.Response
	ID       string `json:"ID"`
	Name     string `json:"name"`
	Age      int    `json:"age" `
	Position string `json:"position"`
}
