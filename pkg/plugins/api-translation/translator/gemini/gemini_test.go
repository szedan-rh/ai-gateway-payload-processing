/*
Copyright 2026 The opendatahub.io Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gemini

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTranslateRequest_BasicChat(t *testing.T) {
	body := map[string]any{
		"model":    "gemini-2.0-flash",
		"messages": []any{map[string]any{"role": "user", "content": "Hi"}},
	}

	translatedBody, headers, headersToRemove, err := NewGeminiTranslator().TranslateRequest(body)
	require.NoError(t, err)
	assert.Nil(t, translatedBody, "Gemini OpenAI-compatible API should not mutate the request body")
	assert.Equal(t, "/v1beta/openai/chat/completions", headers[":path"])
	assert.Len(t, headers, 1)
	assert.Nil(t, headersToRemove)
}

func TestTranslateRequest_PassthroughAllChatParams(t *testing.T) {
	body := map[string]any{
		"model": "gemini-2.0-flash",
		"messages": []any{
			map[string]any{"role": "system", "content": "You are helpful."},
			map[string]any{"role": "user", "content": "Hello"},
		},
		"temperature":       0.7,
		"top_p":             0.9,
		"max_tokens":        1000,
		"stream":            true,
		"stop":              []any{"END"},
		"n":                 1,
		"presence_penalty":  0.5,
		"frequency_penalty": 0.3,
	}

	translatedBody, headers, headersToRemove, err := NewGeminiTranslator().TranslateRequest(body)
	require.NoError(t, err)
	assert.Nil(t, translatedBody, "Gemini OpenAI-compatible API should not mutate the request body")
	assert.Equal(t, "/v1beta/openai/chat/completions", headers[":path"])
	assert.Len(t, headers, 1)
	assert.Nil(t, headersToRemove)
}

func TestTranslateRequest_MissingModel(t *testing.T) {
	body := map[string]any{
		"messages": []any{map[string]any{"role": "user", "content": "Hi"}},
	}

	_, _, _, err := NewGeminiTranslator().TranslateRequest(body)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "model field is required")
}

func TestTranslateRequest_EmptyModel(t *testing.T) {
	body := map[string]any{
		"model":    "",
		"messages": []any{map[string]any{"role": "user", "content": "Hi"}},
	}

	_, _, _, err := NewGeminiTranslator().TranslateRequest(body)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "model field is required")
}

func TestTranslateRequest_EmptyMessages(t *testing.T) {
	body := map[string]any{
		"model":    "gemini-2.0-flash",
		"messages": []any{},
	}

	_, _, _, err := NewGeminiTranslator().TranslateRequest(body)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "messages")
	assert.Contains(t, err.Error(), "BadRequest")
}

func TestTranslateRequest_MissingMessages(t *testing.T) {
	body := map[string]any{
		"model": "gemini-2.0-flash",
	}

	_, _, _, err := NewGeminiTranslator().TranslateRequest(body)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "messages")
}

func TestTranslateResponse_Passthrough(t *testing.T) {
	body := map[string]any{
		"id":      "chatcmpl-abc123",
		"object":  "chat.completion",
		"created": 1700000000,
		"model":   "gemini-2.0-flash",
		"choices": []any{
			map[string]any{
				"index": 0,
				"message": map[string]any{
					"role":    "assistant",
					"content": "The answer is 4.",
				},
				"finish_reason": "stop",
			},
		},
		"usage": map[string]any{
			"prompt_tokens":     10,
			"completion_tokens": 5,
			"total_tokens":      15,
		},
	}

	translatedBody, err := NewGeminiTranslator().TranslateResponse(body, "gemini-2.0-flash")
	require.NoError(t, err)
	assert.Nil(t, translatedBody, "Gemini OpenAI-compatible response should not be mutated")
}

func TestTranslateResponse_NoError(t *testing.T) {
	body := map[string]any{
		"error": map[string]any{
			"message": "Model not found",
			"type":    "invalid_request_error",
			"code":    "model_not_found",
		},
	}

	translatedBody, err := NewGeminiTranslator().TranslateResponse(body, "invalid-model")
	require.NoError(t, err)
	assert.Nil(t, translatedBody, "Error responses should pass through unchanged")
}
