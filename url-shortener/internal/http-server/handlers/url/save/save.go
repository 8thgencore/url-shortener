package save

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	utils "url-shortener/internal/lib/helpers"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

// TODO: move to config
const aliasLength = 4

//go:generate go run github.com/vektra/mockery/v2 --name=URLSaver --case=snake
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
	AliasExists(alias string) (bool, error)
	URLExists(urlToCheck string) (bool, error)
	GetAliasByURL(urlToFind string) (string, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			// Такую ошибку встретим, если получили запрос с пустым телом.
			// Обработаем её отдельно
			log.Error("request body is empty")

			render.JSON(w, r, response.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request", sl.Err(err))

			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("request validation failed", sl.Err(err))

			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" {
			// Generate a random alias until a unique one is found
			exists, err := urlSaver.URLExists(req.URL)
			if err != nil {
				log.Error("failed to check that URL exists in DB", sl.Err(err))
				render.JSON(w, r, response.Error("failed to check that URL exists in DB"))
				return
			}

			if exists {
				alias, err = urlSaver.GetAliasByURL(req.URL)
				if err != nil {
					log.Error("failed to get alias connected to URL", sl.Err(err))
					render.JSON(w, r, response.Error("failed to get alias connected to URL"))
					return
				}

				responseOK(w, r, alias)
				return
			}

			const maxAttempts = 64 // Maximum number of generation attempts
			exists = true

			for attempt := 1; attempt <= maxAttempts; attempt++ {
				alias = random.NewRandomString(aliasLength)
				exists, err = urlSaver.AliasExists(alias)
				if err != nil {
					log.Error("failed to generate alias", sl.Err(err))
					render.JSON(w, r, response.Error("failed to generate url"))
					return
				}
				if !exists {
					break
				}
			}

			if exists {
				log.Error("The number of attempts to create an alias has been exceeded", sl.Err(err))
				render.JSON(w, r, response.Error("The number of attempts to create an alias has been exceeded. Try again after a while"))
				return
			}
		}

		if !utils.IsValidAlias(alias) {
			log.Info("url alias not valid", slog.String("alias", req.URL))

			render.JSON(w, r, response.Error("url alias not valid"))

			return
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, response.Error("url already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, response.Error("failed to add url"))

			return
		}

		log.Info("url saved", slog.Int64("id", id))

		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: response.Ok(),
		Alias:    alias,
	})
}
