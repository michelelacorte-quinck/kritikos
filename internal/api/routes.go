package api

import (
	"kritikos/internal/api/handlers"
	"kritikos/pkg/httpcore"

	"github.com/go-chi/chi/v5"
)

func ApplyRoutes(router chi.Router, controller *handlers.ApiController) {
	router.Post("/ai/kritikos", httpcore.Handle(controller.GetKritikos))
}
