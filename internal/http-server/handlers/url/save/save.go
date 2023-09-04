package save

import (
	"errors"
	resp "github.com/AlexandrLitkevich/qwery/internal/lib/api/response"
	"github.com/AlexandrLitkevich/qwery/internal/lib/logger/sl"
	"github.com/AlexandrLitkevich/qwery/internal/lib/random"
	"github.com/AlexandrLitkevich/qwery/internal/storage"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

//omitempty если нет ошибки то этого поля не будет в JSON

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

// TODO: maybe move to config if needed
const aliasLength = 6

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

// New 2ым параметром передаем интрефейс
func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New" //operation for message error

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())), //трейснг запросов
		)

		var req Request
		// парсим json в структуру Request
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

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		//TODO: проверить есть ли уже созданный alias
		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))
			render.JSON(w, r, resp.Error("url already exists"))
			return
		}

		if err != nil {
			log.Error("failed to add url", sl.Err(err))
			/*
				Ошибка storage передается только в логи
				так как там может содержаться важная информация(БД и тд)
				На клиент передается только текстовая ошибка
			*/
			render.JSON(w, r, resp.Error("failed to add url"))
			return
		}

		log.Info("url added", slog.Int64("id", id))
		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
