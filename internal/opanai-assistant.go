/*
The OpenAI implementation of the LLM
*/

package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type openAIModel string

var openAIClient http.Client = http.Client{}

const (
	GPT_3_5_turbo openAIModel = "gpt-3.5-turbo"
)
const openAIURL string = "https://api.openai.com/v1/chat/completions"

type openAIAssistant struct {
	APIKey        string
	ModelName     openAIModel
	SystemPrompt  string
	UserRole      string
	AssistantRole string
}

type apiRequestPayload struct {
	Model       openAIModel       `json:"model"`
	Messages    []DialogueElement `json:"messages"`
	Temperature float32           `json:"temperature"`
}

type choice struct {
	Index        int             `json:"index"`
	Message      DialogueElement `json:"message"`
	LogProbs     interface{}     `json:"logprobs"`
	FinishReason string          `json:"finish_reason"`
}

type usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type apiResponsePayload struct {
	ID                string      `json:"id"`
	Object            string      `json:"object"`
	Created           int         `json:"created"`
	Model             string      `json:"model"`
	SystemFingerprint string      `json:"system_fingerprint"`
	Choices           []choice    `json:"choices"`
	LogProbs          interface{} `json:"logprobs"`
	FinishReason      string      `json:"finish_reason"`
	Usage             usage       `json:"usage"`
}

func unmarshalOpenAIResponse(resp *http.Response) DialogueElement {
	var payload apiResponsePayload
	err := json.NewDecoder(resp.Body).Decode(&payload)
	PanicOnErr(err)
	return payload.Choices[0].Message
}

// Actually hit the OpenAI API
func (ass *openAIAssistant) hitLargeLanguageModel(client *http.Client, dialogue []DialogueElement) DialogueElement {
	body, err := json.Marshal(apiRequestPayload{
		Model:       ass.ModelName,
		Messages:    dialogue,
		Temperature: 0.7,
	})
	PanicOnErr(err)
	req, err := http.NewRequest("POST", openAIURL, bytes.NewBuffer(body))
	PanicOnErr(err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", ass.APIKey))
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	PanicOnErr(err)
	InspectAPIResponsePayload(resp)
	defer resp.Body.Close()
	return unmarshalOpenAIResponse(resp)
}

// Constructs the system prompt the way OpenAI expects
func (ass *openAIAssistant) createSystemPrompt() DialogueElement {
	return DialogueElement{Role: "system", Content: ass.SystemPrompt}
}

// takes the dialoge provided and passes it to the OpenAI LLM.
// Returns OpenAI's parsed response as a string.
func (ass *openAIAssistant) Prompt(convo *Conversation) string {
	if ass.SystemPrompt != "" {
		system_message := ass.createSystemPrompt()
		convo.AddSystemPrompt(system_message)
	}
	new_dialogue := ass.hitLargeLanguageModel(&openAIClient, convo.History)
	return new_dialogue.Content
}

func (ass *openAIAssistant) GetUserRole() string {
	return ass.UserRole
}

func (ass *openAIAssistant) GetAssistantRole() string {
	return ass.AssistantRole
}

func NewOpenAIAssistant(api_key string, system_prompt string) Assistant {
	if api_key == "" {
		panic(ManagedError{
			Message: "No OpenAI API key provided."})
	}
	return &openAIAssistant{ModelName: GPT_3_5_turbo, APIKey: api_key, SystemPrompt: system_prompt, UserRole: "user", AssistantRole: "assistant"}
}
