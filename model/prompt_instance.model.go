package model

import (
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/prompts"
)

type PromptInstance struct {
	Prompt llms.PromptValue `json:"prompt"`
}

func NewPromptInstance() *PromptInstance {
	return &PromptInstance{}
}

func (p *PromptInstance) GetPrompt(rawText string) string {

	promptParam := prompts.PromptTemplate{
		Template: `

	Based on the provided content: " {{.rawText}} ", perform the following tasks:

	### Main Objective
	xxx.

	### Processing Instructions

	1. **Contextual Analysis**
	- xxx
	- xxx
	- xxx

	2. **Dialogue Structuring**
	- xxx

	3. **Finalization Criteria**
	- xxx
		* xxx
		* xxx
		* xxx
	- xxx

	4. **Response Standards**
	- xxx
	- xxx
	- xxx

	### Required Output Format
	xxx
    
	Output Format: xxx
	{
		"param1": "[xxx]",
		"param2": "[xxx]",
		"param3": "[xxx]",
		"param4": "[xxx]"
	}
	
	`,
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
