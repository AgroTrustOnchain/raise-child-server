package ai

import (
	"context"
	"log"
	"os"
	"raise-child/constants/env"
	"time"

	"github.com/google/generative-ai-go/genai"
	"golang.org/x/time/rate"
	"google.golang.org/api/option"
)

type geminiProvider struct {
	client  *genai.Client
	model   *genai.GenerativeModel
	limiter *rate.Limiter
}

func initializeGeminiClient(ctx context.Context, errLogger *log.Logger) *geminiProvider {
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv(env.GEMINI_API_KEY)))
	if err != nil {
		errLogger.Println(err.Error())
		return nil
	}

	var model = client.GenerativeModel("gemini-2.5-flash")
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text(_prompt_instruction),
		},
	}

	return &geminiProvider{
		client:  client,
		model:   model,
		limiter: rate.NewLimiter(rate.Every(time.Second*12), 1), // 5 reqs/min
	}
}
