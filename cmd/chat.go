package cmd

import (
	"github.com/andrew-candela/chatbot/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var speakOption bool

func init() {
	rootCMD.AddCommand(chatCommand)
	chatCommand.Flags().BoolVarP(&speakOption, "speak-output", "s", false, "speaks the output of the LLM using GNU say")
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
		ass := internal.NewOpenAIAssistant(viper.GetString("openai_api_key"))
		handlers := []internal.ModelOutputHandler{internal.NewStdOutHandler()}
		if speakOption {
			handlers = append(handlers, internal.NewSpeakHandler())
		}
		llm := internal.LLM{
			Assistant:      ass,
			OutputHandlers: handlers,
			SystemPrompt:   viper.GetString("openai_system_prompt"),
		}
		llm.InputLoop()
	},
}
