package connections

import (
	"context"
	"fmt"

	"github.com/Andhika-GIT/go-message-broker-monorepo/configs"
)

func NewGenkit(ctx context.Context, cfg configs.LLMConfig) (*genkit.Genkit, string, error) {
	switch cfg.Provider {
	case "anthropic":
		g, err := genkit.Init(ctx, genkit.WithPlugins(&anthropic.Anthropic{
			APIKey: cfg.APIKey,
		}))
		return g, "anthropic/" + cfg.Model, err

	case "openai_compatible":
		g, err := genkit.Init(ctx, genkit.WithPlugins(&compat_oai.OpenAICompatible{
			Provider: "myllm",
			APIKey:   cfg.APIKey,
			BaseURL:  cfg.BaseURL,
		}))
		return g, "myllm/" + cfg.Model, err

	default:
		return nil, "", fmt.Errorf("unknown provider: %s", cfg.Provider)
	}
}
