package create

import (
	"context"
	"errors"
	"github.com/AlexandrLitkevich/qwery/internal/http-server/handlers/user"
	resp "github.com/AlexandrLitkevich/qwery/internal/lib/api/response"
	"github.com/AlexandrLitkevich/qwery/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type Provider interface {
	Create(ctx context.Context, cancel context.CancelFunc, log *slog.Logger, request user.Request) (bool, error)
}

func New(ctx context.Context, log *slog.Logger, method Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.create.Create"
		ctw, cancel := context.WithTimeout(ctx, 5*time.Second) //TODO send to config
		defer cancel()

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Debug("Start create request")

		var req user.Request
		// парсим json в структуру из Body in Request
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			// Такую ошибку встретим, если получили запрос с пустым телом.
			// Обработаем её отдельно
			log.Error("request body is empty")
			render.JSON(w, r, resp.Error("empty request"))
			return
		}

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		//TODO: added validator
		res, err := method.Create(ctw, cancel, log, req)
		log.Info("This result create etcd", res)

		if err != nil {
			log.Error("failed to add user", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to add user"))
			return
		}
		log.Info("user created")
		//TODO: fix
		responseOK(w, r, user.Response{
			ID:       "1",
			Name:     "",
			Age:      0,
			Position: "",
		})
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, data user.Response) {
	render.JSON(w, r, user.Response{
		Response: resp.OK(),
		ID:       data.ID,
		Name:     data.Name,
		Age:      data.Age,
		Position: data.Position,
	})
}
