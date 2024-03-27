/*
Implements the LLM using the Claude APIs from Anthropic
https://docs.anthropic.com/claude/reference/getting-started-with-the-api
*/

package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type claudeAIModel string

var claudeAIClient http.Client = http.Client{}

const (
	Haiku claudeAIModel = "claude-3-haiku-20240307"
)
const anthropicURL string = "https://api.anthropic.com/v1/messages"

type anthropicAssistant struct {
	APIKey       string
	ModelName    claudeAIModel
	SystemPrompt string
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
	Usage        anthropicResponseUsage     `json:"usage"`
}

type anthropicAPIRequestPayload struct {
	Model     claudeAIModel     `json:"model"`
	MaxTokens int               `json:"max_tokens"`
	Messages  []DialogueElement `json:"messages"`
	System    string            `json:"system"`
}

type anthropicError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type anthropicErrorResponse struct {
	Type  string         `json:"type"`
	Error anthropicError `json:"error"`
}

func (ass *anthropicAssistant) unmarshalAnthropicAPIResponse(resp *http.Response) DialogueElement {
	var payload anthropicResponsePayload
	if resp.StatusCode/100 == 4 {
		var errorPayload anthropicErrorResponse
		err := json.NewDecoder(resp.Body).Decode(&errorPayload)
		PanicOnErr(err)
		panic(
			ManagedError{
				fmt.Sprintf(
					"Error Code %v: %v\nReason :%v",
					resp.StatusCode,
					errorPayload.Error.Type,
					errorPayload.Error.Message,
				),
			},
		)
	}
	if resp.StatusCode >= 500 {
		panic(ManagedError{fmt.Sprintf("Got error code %v from Anthropic API", resp.StatusCode)})
	}
	err := json.NewDecoder(resp.Body).Decode(&payload)
	PanicOnErr(err)

	return DialogueElement{
		Role:    ass.GetAssistantRole(),
		Content: payload.Content[0].Text,
	}
}

func (ass *anthropicAssistant) hitLargeLanguageModel(client *http.Client, dialogue []DialogueElement) DialogueElement {
	body, err := json.Marshal(anthropicAPIRequestPayload{
		Model:     ass.ModelName,
		MaxTokens: 1024,
		Messages:  dialogue,
		System:    ass.SystemPrompt,
	})
	PanicOnErr(err)
	req, err := http.NewRequest("POST", anthropicURL, bytes.NewBuffer(body))
	PanicOnErr(err)
	req.Header.Add("x-api-key", ass.APIKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("anthropic-version", "2023-06-01")
	resp, err := client.Do(req)
	PanicOnErr(err)
	defer resp.Body.Close()
	return ass.unmarshalAnthropicAPIResponse(resp)

}

func (ass *anthropicAssistant) Prompt(convo *Conversation) string {
	new_dialogue := ass.hitLargeLanguageModel(&claudeAIClient, convo.History)
	return new_dialogue.Content
}

func (ass *anthropicAssistant) GetAssistantRole() string {
	return "assistant"
}

func (ass *anthropicAssistant) GetUserRole() string {
	return "user"
}

func NewAnthropicAssistant(api_key string, model_name string, system_prompt string) Assistant {
	if api_key == "" {
		panic(ManagedError{Message: "No Anthropic API key provided."})
	}
	return &anthropicAssistant{ModelName: Haiku, APIKey: api_key, SystemPrompt: system_prompt}
}
