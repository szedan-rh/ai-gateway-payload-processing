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
	"github.com/opendatahub-io/ai-gateway-payload-processing/pkg/plugins/api-translation/translator"
	errcommon "sigs.k8s.io/gateway-api-inference-extension/pkg/common/error"
)

const (
	// Gemini's OpenAI-compatible endpoint path.
	// Google Gemini serves OpenAI-compatible chat completions at /v1beta/openai/chat/completions
	// rather than the standard /v1/chat/completions.
	geminiPath = "/v1beta/openai/chat/completions"
)

// compile-time interface check
var _ translator.Translator = &GeminiTranslator{}

// NewGeminiTranslator initializes a new GeminiTranslator and returns its pointer.
func NewGeminiTranslator() *GeminiTranslator {
	return &GeminiTranslator{}
}

// GeminiTranslator rewrites the request path to Gemini's OpenAI-compatible endpoint.
// Gemini accepts the same request/response schema as OpenAI, so only path rewriting is needed.
type GeminiTranslator struct{}

// TranslateRequest rewrites the path to target Gemini's OpenAI-compatible endpoint.
// The request body is not mutated since Gemini's OpenAI-compatible API accepts the same schema as OpenAI.
func (t *GeminiTranslator) TranslateRequest(body map[string]any) (translatedBody map[string]any,
	headersToMutate map[string]string, headersToRemove []string, err error) {
	model, _ := body["model"].(string)
	if model == "" {
		return nil, nil, nil, errcommon.Error{Code: errcommon.BadRequest, Msg: "model field is required"}
	}

	messages, _ := body["messages"].([]any)
	if len(messages) == 0 {
		return nil, nil, nil, errcommon.Error{Code: errcommon.BadRequest, Msg: "messages field is required and must not be empty"}
	}

	headersToMutate = map[string]string{
		":path": geminiPath,
	}

	// Return nil body — no body mutation needed, Gemini accepts OpenAI request format as-is.
	return nil, headersToMutate, nil, nil
}

// TranslateResponse is a no-op since Gemini's OpenAI-compatible API returns responses in OpenAI format.
func (t *GeminiTranslator) TranslateResponse(body map[string]any, model string) (translatedBody map[string]any, err error) {
	return nil, nil
}
