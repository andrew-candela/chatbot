package cmd

import (
	"github.com/andrew-candela/chatbot/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCMD.AddCommand(initCommand)
}

var initCommand = &cobra.Command{
	Use:   "init",
	Short: "Create the chatbot config file",
	Long: `
	Generates a ~/chatbot/config.toml file.
	You will have to manually edit it as needed.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.InitChatbot()
	},
}
