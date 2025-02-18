package llm_utils

import (
	"context"
	"fmt"
	"time"
)

type LLM interface {
	Chat(ctx context.Context, chatHistory []ChatMessage) (ChatMessage, error)
}

type LLMUtils struct {
	LLM        LLMModel
	Retries    int
	RetryDelay time.Duration
}

func NewLLMUtils(llm LLMModel, retries int, retryDelay time.Duration) *LLMUtils {
	return &LLMUtils{
		LLM:        llm,
		Retries:    retries,
		RetryDelay: retryDelay,
	}
}

func (u *LLMUtils) withRetries(ctx context.Context, fn func() (ChatMessage, error)) (ChatMessage, error) {
	var err error
	for i := 0; i < u.Retries; i++ {
		res, ferr := fn()
		if ferr == nil {
			return res, nil
		}
		err = ferr
		time.Sleep(u.RetryDelay)
	}
	return ChatMessage{}, fmt.Errorf("after %d attempts: %w", u.Retries, err)
}

func SimpleInstruct(ctx context.Context, llm LLM, prompt string) (string, error) {
	history := []ChatMessage{
		{Role: "user", Content: prompt},
	}
	res, err := llm.Chat(ctx, history)
	return res.Content, err
}

func SimpleSummaryTruncate(ctx context.Context, llm LLM, content string, maxLength int) (string, error) {
	if len(content) > maxLength {
		content = content[:maxLength]
	}

	history := []ChatMessage{
		{Role: "system", Content: "Please provide a summary of the following content"},
		{Role: "user", Content: content},
		{Role: "system", Content: "Now please summarize the content above, don't do anything else"},
	}

	res, err := llm.Chat(ctx, history)
	return res.Content, err
}

func BooleanTwoStep(ctx context.Context, llm LLM, content, instruction string) (bool, error) {
	history := []ChatMessage{
		{Role: "system", Content: "Determine the answer to the question about this content"},
		{Role: "user", Content: content},
		{Role: "system", Content: "First, think through the answer to: " + instruction},
	}

	// First step - reasoning
	res, err := llm.Chat(ctx, history)
	if err != nil {
		return false, err
	}

	// Second step - final answer
	history = append(history, res)
	history = append(history, ChatMessage{
		Role:    "user",
		Content: "Answer only 'yes' or 'no'",
	})

	finalRes, err := llm.Chat(ctx, history)
	if err != nil {
		return false, err
	}

	switch finalRes.Content {
	case "yes":
		return true, nil
	case "no":
		return false, nil
	default:
		return false, fmt.Errorf("invalid response: %s", finalRes.Content)
	}
}
