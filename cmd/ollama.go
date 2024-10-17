package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func addOllama() {
	rootCmd.AddCommand(ollamaCmd)
}

func ollamaFlags() {
	addFlagOllamaUrl(ollamaCmd)
	if err := viper.BindPFlags(ollamaCmd.Flags()); err != nil {
		slog.Error("Cannot bind ollama flags", "err", err)
	}
}

const (
	flagURL = "url"
)

func addFlagOllamaUrl(cmd *cobra.Command) {
	cmd.Flags().String(flagURL, "http://llama-1.its.unibas.ch:11434", "Ollama URL")
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

func getOllamaClient(model string) (*ollama.LLM, error) {
	url := viper.GetString(flagURL)
	slog.Info("connecting to ollama", "model", model, "url", url)
	return ollama.New(
		ollama.WithModel(model),
		ollama.WithServerURL(url),
	)
}

func ollamaChat(ctx context.Context) error {
	llm, err := getOllamaClient("llama3.1")
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
