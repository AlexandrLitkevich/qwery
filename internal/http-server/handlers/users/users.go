package users

import (
	"context"
	resp "github.com/AlexandrLitkevich/qwery/internal/lib/api/response"
	"github.com/AlexandrLitkevich/qwery/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"time"
)

type Response struct {
	resp.Response
	User string
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=GetUser
type UserGetting interface {
	GetUser(ctx context.Context, log *slog.Logger, userId string) (string, error)
}

func New(ctx context.Context, log *slog.Logger, method UserGetting) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.New"

		cth, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())), //трейснг запросов
		)

		userId := chi.URLParam(r, "userId")
		if userId == "" {
			log.Info("userId is empty")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		user, err := method.GetUser(cth, log, userId)
		if err != nil {
			log.Error("fail to get etcd in func:", op, sl.Err(err))
			render.JSON(w, r, resp.Error("invalid response"))
			return
		}
		log.Info("This git user", user)
		responseOK(w, r, user)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, user string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		User:     user,
	})
}
