package models

type AiModel string

type AiResponse struct {
	Content      string `json:"content"`
	InputTokens  int    `json:"inputTokens"`
	OutputTokens int    `json:"outputTokens"`
	TotalTokens  int    `json:"totalTokens"`
}

type KritikosRequest struct {
	SystemPrompt                string  `json:"systemPrompt" validate:"required"`
	Prompt                      string  `json:"prompt" validate:"required"`
	BaseModel                   string  `json:"baseModel" validate:"required"`
	EvaluationModel             string  `json:"evaluationModel" validate:"required"`
	BaseModelTemperature        float32 `json:"baseModelTemperature"`
	EvalidationModelTemperature float32 `json:"evaluationModelTemperature"`
	MaxRetries                  int     `json:"maxRetries" validate:"required"`
}

type EvaluationInput struct {
	UserSystemPrompt string `json:"userSystemPrompt" validate:"required"`
	UserPrompt       string `json:"userPrompt" validate:"required"`
	DraftAnswer      string `json:"draftAnswer" validate:"required"`
	ModelEvaluation  string `json:"modelEvaluation"`
}

type EvaluationResult struct {
	Scores struct {
		Relevance    int `json:"relevance" validate:"required"`
		Correctness  int `json:"correctness" validate:"required"`
		Completeness int `json:"completeness" validate:"required"`
		Clarity      int `json:"clarity" validate:"required"`
		Style        int `json:"style" validate:"required"`
	} `json:"scores"`
	Strengths        []string `json:"strengths" validate:"required"`
	Weaknesses       []string `json:"weaknesses" validate:"required"`
	ActionableAdvice []string `json:"actionableAdvice" validate:"required"`
	ImprovedAnswer   string   `json:"improvedAnswer" validate:"required"`
}
