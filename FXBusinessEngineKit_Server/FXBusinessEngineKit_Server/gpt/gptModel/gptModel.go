package gptModel

type UserChatModel struct {
	Content string `json:"content"`
}

type ChatModelUsageModel struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptToken      int `json:"prompt_token"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatModelChoicesMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatModelChoices struct {
	Index        int                     `json:"index"`
	FinishReason string                  `json:"finish_reason"`
	Message      ChatModelChoicesMessage `json:"message"`
}

type ChatModel struct {
	Id      string              `json:"id"`
	Object  string              `json:"object"`
	Created int                 `json:"created"`
	Model   string              `json:"model"`
	Choices []ChatModelChoices  `json:"choices"`
	Usage   ChatModelUsageModel `json:"usage"`
}

type ChatRequestModelMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequestModel struct {
	Messages         []ChatRequestModelMessage `json:"messages"`
	Temperature      float32                   `json:"temperature"`
	TopP             float32                   `json:"top_p"`
	FrequencyPenalty float32                   `json:"frequency_penalty"`
	PresencePenalty  float32                   `json:"presence_penalty"`
	MaxTokens        int                       `json:"max_tokens"`
}
