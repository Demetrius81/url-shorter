package remove

import (
	"errors"
	"log/slog"
	"net/http"

	httpresponse "url-shorter/internal/lib/api/http-response"
	resp "url-shorter/internal/lib/api/response"
	"url-shorter/internal/lib/logger/sl"
	"url-shorter/internal/storage"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type URLRemover interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlRemover URLRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		err := urlRemover.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)

			render.JSON(w, r, resp.Error("not found"))

			return
		}
		if err != nil {
			log.Error("failed to delete url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("url is deleted", slog.String("alias", alias))

		responseOk(w, r, alias)

	}

}

func responseOk(w http.ResponseWriter, r *http.Request, alias string) {

	render.JSON(w, r, httpresponse.Response{Response: resp.OK(), Alias: alias})
}
