package cmd

import (
	"github.com/andrew-candela/chatbot/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCMD.AddCommand(chatCommand)
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
		output := internal.NewStdOutHandler()
		llm := internal.LLM{
			Assistant:     ass,
			OutputHandler: output,
			SystemPrompt:  viper.GetString("openai_system_prompt"),
		}
		llm.InputLoop()
	},
}
