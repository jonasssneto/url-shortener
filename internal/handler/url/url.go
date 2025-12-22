package url_handler

import (
	"encoding/json"
	url_dto "main/internal/dto/url"
	"main/internal/metrics"
	usecase "main/internal/use-case/url"
	errorhandler "main/internal/utils/error"
	"main/pkg/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type URLHandler struct {
	Usecase *usecase.URLUseCase
	Logger  *logger.Logger
}

func New(usecase *usecase.URLUseCase) *URLHandler {
	return &URLHandler{
		Usecase: usecase,
		Logger:  logger.New("url-handler"),
	}
}

func (u *URLHandler) Create(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("url-create-handler")
	ctx, span := tracer.Start(r.Context(), "create-handler")
	defer span.End()

	var dto url_dto.CreateURLDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to decode request body")
		errorhandler.JsonError(
			w,
			"Erro ao decodificar o corpo da requisição",
			http.StatusBadRequest,
			map[string]interface{}{
				"error": err.Error(),
			},
		)
		metrics.UrlsCreated.WithLabelValues("error").Inc()
		return
	}

	span.SetAttributes(
		attribute.String("url.original", dto.OriginalURL),
		attribute.String("url.slug", dto.Slug),
	)

	u.Logger.Debugf("Received Create request: %+v\n", dto)

	err = u.Usecase.Create(ctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create URL")

		errorhandler.JsonError(
			w,
			"Erro ao criar URL",
			http.StatusBadRequest,
			map[string]interface{}{
				"error": err.Error(),
				"input": dto,
			},
		)

		metrics.UrlsCreated.WithLabelValues("error").Inc()
		return
	}

	metrics.UrlsCreated.WithLabelValues("success").Inc()
	span.SetAttributes(
		attribute.Bool("url.created", true),
	)

	w.WriteHeader(http.StatusCreated)
}

func (u *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("url-redirect-handler")
	ctx, span := tracer.Start(r.Context(), "redirect-handler")
	defer span.End()

	slug := chi.URLParam(r, "slug")
	if slug == "" {
		errorhandler.JsonError(
			w,
			"Slug não informado",
			http.StatusNotFound,
			map[string]interface{}{
				"input": r.URL.String(),
			},
		)
		metrics.UrlsRedirected.WithLabelValues("error").Inc()
		return
	}

	span.SetAttributes(
		attribute.String("url.slug", slug),
	)

	dto := url_dto.RedirectURLDTO{
		Slug: slug,
	}

	u.Logger.Debugf("Received Redirect request for slug: %s", slug)

	url, err := u.Usecase.Redirect(ctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to redirect")
		errorhandler.JsonError(
			w,
			"Erro ao redirecionar",
			http.StatusBadRequest,
			map[string]interface{}{
				"error": err.Error(),
				"input": dto,
			},
		)
		metrics.UrlsRedirected.WithLabelValues("error").Inc()
		return
	}

	metrics.UrlsRedirected.WithLabelValues("success").Inc()
	span.SetAttributes(
		attribute.Bool("redirect.success", true),
		attribute.String("redirect.target", url.OriginalURL),
	)

	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}
