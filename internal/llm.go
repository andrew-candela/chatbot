package internal

import (
	"fmt"

	"github.com/chzyer/readline"
)

// Handles communicating the model output to the user.
// Can be something that prints to stdout or actually speaks.
type ModelOutputHandler interface {
	Output(string)
}

// The interface for the LLM
type Assistant interface {
	Prompt(*Conversation) string
}

type DialogueElement struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Conversation struct {
	History []DialogueElement
}

func (convo *Conversation) ChatHistory() []DialogueElement {
	return convo.History
}

func (convo *Conversation) AddDialogue(element DialogueElement) {
	convo.History = append(convo.History, element)
}

type LLM struct {
	Assistant     Assistant
	OutputHandler ModelOutputHandler
	SystemPrompt  string
	Conversation  Conversation
}

// Adds a beginning prompt to the conversation before taking
// input from the user.
func addSystemPrompt(prompt string, conversation *Conversation) {
	new_element := DialogueElement{Role: "system", Content: prompt}
	conversation.AddDialogue(new_element)
}

/*
Collects input from the user and sends it to the LLM,
then handles the output.
*/
func (llm *LLM) InputLoop() {
	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	if llm.SystemPrompt != "" {
		addSystemPrompt(llm.SystemPrompt, &llm.Conversation)
	}

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF
			break
		}
		if line != "" {
			fmt.Println()
			input := DialogueElement{
				Role:    "user",
				Content: line,
			}
			llm.Conversation.AddDialogue(input)
			response := llm.Assistant.Prompt(&llm.Conversation)
			output := DialogueElement{
				Role:    "assistant",
				Content: response,
			}
			llm.Conversation.AddDialogue(output)
			llm.OutputHandler.Output(output.Content)
		}
	}
}
