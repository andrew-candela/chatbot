/*
Implements the LLM using the Claude APIs from Anthropic
https://docs.anthropic.com/claude/reference/getting-started-with-the-api
*/

package internal

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type claudeAIModel string

var claudeAIClient http.Client = http.Client{}

const (
	ATOM claudeAIModel = "atom" // or something
)
const anthropicURL string = "https://api.anthropic.com/v1/messages"

type anthropicAssistant struct {
	APIKey    string
	ModelName claudeAIModel
}

type anthropicResponseContent struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type anthropicResponseUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type anthropicResponsePayload struct {
	Content      []anthropicResponseContent `json:"content"`
	ID           string                     `json:"id"`
	Model        string                     `json:"model"`
	Role         string                     `json:"role"`
	StopReason   string                     `json:"stop_reason"`
	StopSequence string                     `json:"stop_sequence"`
	Type         string                     `json:"type"`
}

func unmarshalAnthropicAPIResponse(resp *http.Response) DialogueElement {
	var payload anthropicResponseContent
	err := json.NewDecoder(resp.Body).Decode(&payload)
	PanicOnErr(err)
	return DialogueElement{
		Role: "assistant"
		Content: ,
		.Content[0]
}

func (ass *anthropicAssistant) hitLargeLanguageModel(client *http.Client, dialogue []DialogueElement) DialogueElement {
	body, err := json.Marshal(apiRequestPayload{
		Model:       ass.ModelName,
		Messages:    dialogue,
		Temperature: 0.7,
	})
	PanicOnErr(err)
	req, err := http.NewRequest("POST", openAIURL, bytes.NewBuffer(body))
	PanicOnErr(err)
	req.Header.Add("x-api-key", ass.APIKey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	PanicOnErr(err)
	defer resp.Body.Close()
	PanicOnErr(err)
	return unmarshalOpenAIResponse(resp)

}

func (ass *anthropicAssistant) Prompt(convo *Conversation) string {
	new_dialogue := ass.hitLargeLanguageModel(&openAIClient, convo.History)
	convo.AddDialogue(new_dialogue)
	return convo.ChatHistory()[len(convo.ChatHistory())-1].Content
}

func NewAnthropicAssistant(api_key string) Assistant {
	if api_key == "" {
		panic(ManagedError{Message: "No ANthropic API key provided."})
	}
	return &anthropicAssistant{ModelName: ATOM, APIKey: api_key}
}
