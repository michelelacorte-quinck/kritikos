package handlers

import (
	"kritikos/pkg/httpcore"
	"kritikos/pkg/models"
	"net/http"
)

func (c *ApiController) GetKritikos(w http.ResponseWriter, r *http.Request) (any, int) {
	ctx := r.Context()

	kritikosRequest, err := httpcore.DecodeBody[models.KritikosRequest](w, r)
	if err != nil {
		return httpcore.ErrBadRequest.With(err), http.StatusInternalServerError
	}
	aiResponse, err := c.ai.GenerateCompletion(ctx, kritikosRequest)
	if err != nil {
		return httpcore.ErrUnkownInternal.With(err), http.StatusInternalServerError
	}

	return aiResponse, http.StatusOK
}
