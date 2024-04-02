package cmd

import (
	"github.com/andrew-candela/chatbot/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var left_system_prompt string
var right_system_prompt string
var initial_prompt string
var left_claude bool
var right_claude bool

func init() {
	rootCMD.AddCommand(conversationCommand)
	conversationCommand.Flags().StringVar(&left_system_prompt, "left-prompt", "", "the system prompt to give to the left assistant")
	conversationCommand.Flags().StringVar(&right_system_prompt, "right-prompt", "", "the system prompt to give to the right assistant")
	conversationCommand.Flags().StringVar(&initial_prompt, "initial-prompt", "hi", "the beginning of the conversation")
	conversationCommand.Flags().BoolVar(&left_claude, "left-claude", false, "use claude for the left assistant")
	conversationCommand.Flags().BoolVar(&right_claude, "right-claude", false, "use claude for the right assistant")
}

var conversationCommand = &cobra.Command{
	Use:   "converse",
	Short: "Have two LLMs talk to each other",
	Long: `
	Light your money on fire and watch as two Large
	Language Models speak nonsense to each other.
	Truly, we live in a wonderous age!
	`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.ParseConfigWithViper()
		var left_ass internal.Assistant
		var right_ass internal.Assistant
		if left_claude {
			left_ass = internal.NewAnthropicAssistant(
				viper.GetString("anthropic_api_key"),
				internal.GetViperValueWithDefault("anthropic_model_name", string(internal.Haiku)),
				left_system_prompt,
			)
		} else {
			left_ass = internal.NewOpenAIAssistant(
				viper.GetString("openai_api_key"),
				left_system_prompt,
			)
		}
		if right_claude {
			right_ass = internal.NewAnthropicAssistant(
				viper.GetString("anthropic_api_key"),
				internal.GetViperValueWithDefault("anthropic_model_name", string(internal.Haiku)),
				right_system_prompt,
			)
		} else {
			right_ass = internal.NewOpenAIAssistant(
				viper.GetString("openai_api_key"),
				right_system_prompt,
			)
		}
		left_handlers := []internal.ModelOutputHandler{internal.NewStdOutHandler()}
		right_handlers := []internal.ModelOutputHandler{internal.NewStdOutHandler()}
		leftLLM := internal.LLM{
			Assistant:      left_ass,
			OutputHandlers: left_handlers,
		}
		rightLLM := internal.LLM{
			Assistant:      right_ass,
			OutputHandlers: right_handlers,
		}
		internal.Singularity(&leftLLM, &rightLLM, initial_prompt)
	},
}
