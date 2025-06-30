package handlers

import (
	"kritikos/pkg/aiclient"
)

type ApiController struct {
	ai *aiclient.AiClient
}

func NewApiController(ai *aiclient.AiClient) *ApiController {
	return &ApiController{
		ai: ai,
	}
}
