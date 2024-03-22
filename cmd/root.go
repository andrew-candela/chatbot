package cmd

import (
	"github.com/andrew-candela/chatbot/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCMD = &cobra.Command{
	Use:   "chatbot",
	Short: "talk to a chat assistant",
	Long: `
	Chatbot - hit a LLM with some input and see the results

	Maybe one day I'll host my own LLM for this.`,
}

func Execute() {
	defer internal.CatchPanicAndExit()
	cobra.OnInitialize(initConfig)
	rootCMD.Execute()
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("/etc/chatbot")
	viper.AddConfigPath("$HOME/.chatbot")
	viper.AddConfigPath(".")
}
