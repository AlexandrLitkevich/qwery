package user

import (
	"errors"
	resp "github.com/AlexandrLitkevich/qwery/internal/lib/api/response"
	"github.com/AlexandrLitkevich/qwery/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
)

func Create(log *slog.Logger, UserCRUD Crud) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.createdUser.save.Create" //operation for message error

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())), //трейснг запросов
		)

		log.Debug("Start create request")

		var req Request
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

		if err := validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr) //TODO what error.As ???
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}
		log.Info("This req", req)

		createdUser, err := UserCRUD.CreateUser(req)

		if err != nil {
			log.Error("failed to add user", sl.Err(err))
			/*
				Ошибка storage передается только в логи
				так как там может содержаться важная информация(БД и тд)
				На клиент передается только текстовая ошибка
			*/
			render.JSON(w, r, resp.Error("failed to add user"))
			return
		}

		log.Info("createdUser added", slog.String("createdUser id", createdUser.ID))
		responseOK(w, r, *createdUser)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, user User) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		ID:       user.ID,
		Name:     user.Name,
		Age:      user.Age,
		Position: user.Position,
	})
}
