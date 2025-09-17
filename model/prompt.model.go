package model

import (
	_ "embed"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/prompts"
)

//go:embed example.prompt.md
var prompt_example string

type PromptOCRInstance struct {
	Prompt llms.PromptValue `json:"prompt"`
}

func NewPromptOCRInstance() *PromptOCRInstance {
	return &PromptOCRInstance{}
}

func (p *PromptOCRInstance) GetPrompt(rawText string) string {

	promptParam := prompts.PromptTemplate{
		Template: prompt_example,
		InputVariables: []string{"rawText"},
		TemplateFormat: prompts.TemplateFormatGoTemplate,
	}

	result, err := promptParam.Format(map[string]any{
		"rawText": rawText,
	})
	if err != nil {
		log.Println("Error formatting prompt: ", err)
	}

	return result
}
