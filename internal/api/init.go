package api

import (
	"context"
	"encoding/base64"
	"kritikos/internal/api/handlers"
	"kritikos/pkg/aiclient"
	"kritikos/pkg/httpcore"
	"kritikos/pkg/models"
	"kritikos/pkg/util"
	"net/http"

	"github.com/rs/zerolog/log"
)

func InitService() http.Handler {
	util.InitLogger()

	router := httpcore.NewRouter()
	env := InitEnv()
	ctx := context.Background()

	gCloudCreds, err := base64.StdEncoding.DecodeString(env.GCloudBase64Creds)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to decode base64 Google cloud credentials")
	}

	ai, err := aiclient.NewAiClient(ctx, models.GCloudCredentials{
		ProjectID: env.GCloudProjectID,
		Region:    env.GCloudRegion,
		CredsJson: gCloudCreds,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize AI client")
	}

	controller := handlers.NewApiController(ai)

	ApplyRoutes(router, controller)

	return router
}
