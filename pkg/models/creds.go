package models

type GCloudCredentials struct {
	ProjectID string
	Region    string
	CorpusID  string
	CredsJson []byte
}
