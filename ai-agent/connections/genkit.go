package connections

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/anthropic"
	"github.com/firebase/genkit/go/plugins/compat_oai"
	"github.com/firebase/genkit/go/plugins/googlegenai"

	"github.com/Andhika-GIT/go-message-broker-monorepo/configs"
)

func NewGenkit(ctx context.Context, cfg configs.LLMConfig) (g *genkit.Genkit, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("genkit init failed: %v", r)
		}
	}()

	switch cfg.Provider {
	case "anthropic":
		g = genkit.Init(ctx,
			genkit.WithPlugins(&anthropic.Anthropic{APIKey: cfg.APIKey}),
			genkit.WithDefaultModel("anthropic/"+cfg.Model),
		)
		return g, nil

	case "openai_compatible":
		g = genkit.Init(ctx,
			genkit.WithPlugins(&compat_oai.OpenAICompatible{
				Provider: "myllm",
				APIKey:   cfg.APIKey,
				BaseURL:  cfg.BaseURL,
			}),
			genkit.WithDefaultModel("myllm/"+cfg.Model),
		)
		return g, nil

	case "gemini":
		g = genkit.Init(ctx,
			genkit.WithPlugins(&googlegenai.GoogleAI{APIKey: cfg.APIKey}),
			genkit.WithDefaultModel("googleai/"+cfg.Model),
		)
		return g, nil

	default:
		return nil, fmt.Errorf("unknown provider: %s", cfg.Provider)
	}
}
