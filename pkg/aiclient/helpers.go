package aiclient

import (
	"encoding/json"
	"kritikos/pkg/models"
	"log"

	"google.golang.org/genai"
)

func buildEvaluationSchema() *genai.Schema {
	return &genai.Schema{
		Type:  genai.TypeObject,
		Title: "EvaluationResult",
		Properties: map[string]*genai.Schema{
			"scores":           schemaScores(),
			"strengths":        arrayOfStringSchema(),
			"weaknesses":       arrayOfStringSchema(),
			"actionableAdvice": arrayOfStringSchema(),
			"improvedAnswer":   {Type: genai.TypeString},
		},
		Required: []string{"scores", "strengths", "weaknesses", "actionableAdvice", "improvedAnswer"},
	}
}

func schemaScores() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"relevance":    scoreField(),
			"correctness":  scoreField(),
			"completeness": scoreField(),
			"clarity":      scoreField(),
			"style":        scoreField(),
		},
		Required: []string{"relevance", "correctness", "completeness", "clarity", "style"},
	}
}

func arrayOfStringSchema() *genai.Schema {
	return &genai.Schema{
		Type:  genai.TypeArray,
		Items: &genai.Schema{Type: genai.TypeString},
	}
}

func scoreField() *genai.Schema {
	return &genai.Schema{
		Type:    genai.TypeInteger,
		Minimum: floatPtr(1),
		Maximum: floatPtr(5),
	}
}

func floatPtr(f float64) *float64 {
	return &f
}

func toJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		log.Printf("marshal error: %v", err)
		return ""
	}
	return string(b)
}

func concatenateParts(parts []*genai.Part) string {
	var content string
	for _, p := range parts {
		if p.Text != "" {
			content += p.Text
		}
	}
	return content
}

func safeToken(v int32) int {
	return int(v)
}

func buildEvaluationInput(req models.KritikosRequest, answer, eval string) models.EvaluationInput {
	return models.EvaluationInput{
		UserSystemPrompt: req.SystemPrompt,
		UserPrompt:       req.Prompt,
		DraftAnswer:      answer,
		ModelEvaluation:  eval,
	}
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
