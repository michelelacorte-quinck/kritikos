package aiclient

import (
	"context"
	"encoding/json"
	"fmt"
	"kritikos/pkg/aiclient/prompts"
	"kritikos/pkg/models"
	"log"

	"cloud.google.com/go/auth/credentials"
	"google.golang.org/genai"
)

type AiClient struct {
	gemini *Gemini
}

type Gemini struct {
	client   *genai.Client
	corpusId string
}

const (
	minTemperature        = 0
	temperatureStep       = 0.1
	defaultCandidateCount = 1
)

func NewAiClient(ctx context.Context, gcloudCreds models.GCloudCredentials) (*AiClient, error) {
	creds, err := credentials.DetectDefault(&credentials.DetectOptions{
		Scopes:          []string{"https://www.googleapis.com/auth/cloud-platform"},
		CredentialsJSON: gcloudCreds.CredsJson,
	})
	if err != nil {
		return nil, err
	}

	gemini, err := genai.NewClient(ctx, &genai.ClientConfig{
		Project:     gcloudCreds.ProjectID,
		Location:    gcloudCreds.Region,
		Credentials: creds,
		Backend:     genai.BackendVertexAI,
	})
	if err != nil {
		return nil, err
	}

	return &AiClient{gemini: &Gemini{client: gemini, corpusId: gcloudCreds.CorpusID}}, nil
}

func (a *AiClient) GenerateCompletion(ctx context.Context, req models.KritikosRequest) (models.EvaluationResult, error) {
	res, err := a.generateBaseCompletion(ctx, req)
	if err != nil {
		return models.EvaluationResult{}, err
	}
	return a.performSelfCritique(ctx, req, res)
}

func (a *AiClient) generateBaseCompletion(ctx context.Context, req models.KritikosRequest) (models.AiResponse, error) {
	temp := DefaultBaseModelTemperature
	if req.BaseModelTemperature > 0 {
		temp = req.BaseModelTemperature
	}

	res, err := a.gemini.client.Models.GenerateContent(ctx, req.BaseModel, genai.Text(req.Prompt), &genai.GenerateContentConfig{
		Temperature: &temp,
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: req.SystemPrompt}},
		},
		ResponseMIMEType: "text/plain",
		CandidateCount:   defaultCandidateCount,
	})
	if err != nil {
		return models.AiResponse{}, err
	}

	if len(res.Candidates) == 0 || res.Candidates[0].Content == nil || len(res.Candidates[0].Content.Parts) == 0 {
		return models.AiResponse{}, ErrAiResponseEmpty
	}

	return models.AiResponse{
		Content:      concatenateParts(res.Candidates[0].Content.Parts),
		InputTokens:  safeToken(res.UsageMetadata.PromptTokenCount),
		OutputTokens: safeToken(res.UsageMetadata.CandidatesTokenCount),
		TotalTokens:  safeToken(res.UsageMetadata.TotalTokenCount),
	}, nil
}

func (a *AiClient) performSelfCritique(ctx context.Context, req models.KritikosRequest, aiResp models.AiResponse) (models.EvaluationResult, error) {
	temp := DefaultEvaluationModelTemperature
	maxRetries := DefaultMaxRetries
	if req.MaxRetries > 0 {
		maxRetries = req.MaxRetries
	}

	schema := buildEvaluationSchema()
	prompt := buildEvaluationInput(req, aiResp.Content, "")

	for attempt := 1; attempt <= maxRetries; attempt++ {
		res, err := a.gemini.client.Models.GenerateContent(ctx, req.EvaluationModel, genai.Text(toJSON(prompt)), &genai.GenerateContentConfig{
			Temperature: &temp,
			SystemInstruction: &genai.Content{
				Parts: []*genai.Part{{Text: prompts.KritikosSystemPrompt}},
			},
			ResponseMIMEType: "application/json",
			ResponseSchema:   schema,
			CandidateCount:   defaultCandidateCount,
		})
		if err != nil {
			return models.EvaluationResult{}, fmt.Errorf("generate (attempt %d): %w", attempt, err)
		}

		if len(res.Candidates) == 0 || res.Candidates[0].Content == nil || len(res.Candidates[0].Content.Parts) == 0 {
			return models.EvaluationResult{}, ErrAiResponseEmpty
		}

		var eval models.EvaluationResult
		if err := json.Unmarshal([]byte(concatenateParts(res.Candidates[0].Content.Parts)), &eval); err != nil {
			return models.EvaluationResult{}, fmt.Errorf("json unmarshal (attempt %d): %w", attempt, err)
		}

		if attempt == maxRetries || len(eval.ActionableAdvice) == 0 {
			log.Printf("Self-critique loop: attempt %d succeeded with scores: %+v", attempt, eval.Scores)
			return eval, nil
		}

		log.Printf("Self-critique loop: attempt %d failed with scores: %+v", attempt, eval.Scores)
		temp = float32(max(float64(temp)-temperatureStep, minTemperature))
		prompt = buildEvaluationInput(req, eval.ImprovedAnswer, toJSON(eval))
	}

	return models.EvaluationResult{}, ErrAiResponseNotApproved
}
