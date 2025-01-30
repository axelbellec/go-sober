package llm

import (
	"context"
	"go-sober/platform"
	"log"

	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

// RunLLM executes the LLM generation with the given prompt
func Call(systemPrompt string, userPrompt string) (string, error) {
	// Load the Groq API key from the .env file if you use it
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	llm, err := openai.New(
		openai.WithModel(platform.AppConfig.LLM.Groq.Model),
		openai.WithBaseURL(platform.AppConfig.LLM.Groq.BaseURL),
		openai.WithToken(platform.AppConfig.LLM.Groq.APIKey),
		openai.WithResponseFormat(openai.ResponseFormatJSON),
	)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	completion, err := llm.GenerateContent(ctx, []llms.MessageContent{
		{
			Role:  llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{llms.TextContent{Text: systemPrompt}},
		},
		{
			Role:  llms.ChatMessageTypeGeneric,
			Parts: []llms.ContentPart{llms.TextContent{Text: userPrompt}},
		},
	},
		llms.WithTemperature(0.1),
		llms.WithTopP(1),
		llms.WithMaxTokens(512),
		llms.WithJSONMode(),
	)

	if err != nil {
		log.Fatal(err)
	}

	return completion.Choices[0].Content, nil
}
