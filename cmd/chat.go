package cmd

import (
	"github.com/andrew-candela/chatbot/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var speakOption bool
var useClaude bool
var cli_system_prompt string

func init() {
	rootCMD.AddCommand(chatCommand)
	chatCommand.Flags().BoolVarP(&speakOption, "speak-output", "s", false, "speaks the output of the LLM using GNU say")
	chatCommand.Flags().BoolVarP(&useClaude, "use-claude", "c", false, "uses Anthropic models instead of OpenAI")
	chatCommand.Flags().StringVarP(&cli_system_prompt, "prompt", "p", "", "the system prompt to use for the conversation")
}

var chatCommand = &cobra.Command{
	Use:   "chat",
	Short: "Chat with a LLM",
	Long: `
	Chats with a LLM.
	Stores your conversation history and submits it
	to the LLM on each call.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.ParseConfigWithViper()
		var ass internal.Assistant
		config_system_prompt := viper.GetString("system_prompt")
		system_prompt := internal.Coalesce(cli_system_prompt, config_system_prompt)
		if useClaude {
			ass = internal.NewAnthropicAssistant(
				viper.GetString("anthropic_api_key"),
				internal.GetViperValueWithDefault("anthropic_model_name", string(internal.Haiku)),
				system_prompt,
			)
		} else {
			ass = internal.NewOpenAIAssistant(
				viper.GetString("openai_api_key"),
				system_prompt,
			)
		}
		handlers := []internal.ModelOutputHandler{internal.NewStdOutHandler()}
		if speakOption {
			handlers = append(handlers, internal.NewSpeakHandler())
		}
		llm := internal.LLM{
			Assistant:      ass,
			OutputHandlers: handlers,
		}
		llm.InputLoop()
	},
}
