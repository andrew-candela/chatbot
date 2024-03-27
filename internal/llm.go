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
	GetUserRole() string
	GetAssistantRole() string
}

type DialogueElement struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Conversation struct {
	History         []DialogueElement
	HasSystemPrompt bool
}

func (convo *Conversation) ChatHistory() []DialogueElement {
	return convo.History
}

func (convo *Conversation) AddDialogue(element DialogueElement) {
	convo.History = append(convo.History, element)
}

// Places a system prompt message in the beginning of the conversation.
// If the conversation already has a system prompt, then it
// will replace the first dialogue element with the new system prompt.
func (convo *Conversation) AddSystemPrompt(element DialogueElement) {
	// replace the first element of the conversation with
	// the given system prompt
	if convo.HasSystemPrompt {
		convo.History[0] = element
		return
	}
	// no system prompt exists, so put one into the beginning of the convo
	convo.History = Prepend(convo.History, element)
	convo.HasSystemPrompt = true
}

type LLM struct {
	Assistant      Assistant
	OutputHandlers []ModelOutputHandler
	Conversation   Conversation
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

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF
			break
		}
		if line != "" {
			fmt.Println()
			input := DialogueElement{
				Role:    llm.Assistant.GetUserRole(),
				Content: line,
			}
			llm.Conversation.AddDialogue(input)
			response := llm.Assistant.Prompt(&llm.Conversation)
			output := DialogueElement{
				Role:    llm.Assistant.GetAssistantRole(),
				Content: response,
			}
			llm.Conversation.AddDialogue(output)
			for _, handler := range llm.OutputHandlers {
				handler.Output(output.Content)
			}
		}
	}
}
