package image

import (
	"errors"
	resp "github.com/AlexandrLitkevich/qwery/internal/lib/api/response"
	"github.com/AlexandrLitkevich/qwery/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
)

type Request struct {
	Name string
	Desc string
}
type Response struct {
	resp.Response
	status bool
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=LoadProvider
type LoadProvider interface {
	LoadImage() (bool, error)
}

func New(log *slog.Logger, method LoadProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.image.LoadImage"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())), //трейснг запросов
		)

		log.Debug("Start request", op)

		//var req Request
		_, err := io.ReadAll(r.Body)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.JSON(w, r, resp.Error("empty request"))
			return
		}

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		st, err := method.LoadImage()

		render.JSON(w, r, Response{
			Response: resp.OK(),
			status:   st,
		})
	}
}
