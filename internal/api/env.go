package api

import "kritikos/pkg/util"

type EnvConfig struct {
	GCloudBase64Creds string
	GCloudProjectID   string
	GCloudRegion      string
}

var envVarMappings = util.EnvMapping{
	"GOOGLE_CLOUD_BASE64_CREDS": &env.GCloudBase64Creds,
	"GOOGLE_CLOUD_PROJECT_ID":   &env.GCloudProjectID,
	"GOOGLE_CLOUD_REGION":       &env.GCloudRegion,
}

var env = &EnvConfig{}

func InitEnv() *EnvConfig {
	for key, goVar := range envVarMappings {
		*goVar = util.GetEnvOrPanic(key)
	}

	return env
}
