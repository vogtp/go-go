package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/vogtp/rag/pkg/cfg"
)

func addOllama() {
	rootCmd.AddCommand(ollamaCmd)
}

var ollamaCmd = &cobra.Command{
	Use:     "ollama",
	Short:   "Run ollama stuff",
	Aliases: []string{"o"},
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return ollamaChat(cmd.Context())
	},
}

func getOllamaClient(ctx context.Context, model string) (*ollama.LLM, error) {
	url := cfg.GetOllamaHost(ctx)
	slog.Info("connecting to ollama", "model", model, "url", url)
	return ollama.New(
		ollama.WithModel(model),
		ollama.WithServerURL(url),
	)
}

func ollamaChat(ctx context.Context) error {
	llm, err := getOllamaClient(ctx, viper.GetString(cfg.ModelDefault))
	if err != nil {
		return fmt.Errorf("cannot create ollama connection: %w", err)
	}

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "You are a company branding design wizard."),
		llms.TextParts(llms.ChatMessageTypeHuman, "What would be a good company name for a comapny that produces Go-backed LLM tools?"),
	}
	completion, err := llm.GenerateContent(ctx, content, llms.WithTemperature(0.001), llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		fmt.Print(string(chunk))
		return nil
	}))
	if err != nil {
		return fmt.Errorf("cannot generate content: %w", err)
	}
	_ = completion
	return nil
}
