package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type Client struct {
	apiKey string
	model  string
	http   *http.Client
}

func NewClient(apiKey, model string) *Client {
	if model == "" {
		model = "gemini-2.0-flash"
	}
	return &Client{apiKey: apiKey, model: model, http: &http.Client{Timeout: 20 * time.Second}}
}

type CuriosityResult struct {
	Score     int    `json:"curiosity_score"`
	Reasoning string `json:"reasoning"`
}

// ScoreCuriosity calls the real Gemini API when an API key is configured.
// Without a key, it falls back to a deterministic heuristic so the full
// scoring pipeline still works offline/without credentials.
func (c *Client) ScoreCuriosity(ctx context.Context, textContent []string, primaryEmotion string) (*CuriosityResult, error) {
	if c.apiKey == "" {
		return heuristicCuriosity(textContent, primaryEmotion), nil
	}

	prompt := fmt.Sprintf(`You are a YouTube CTR optimization expert. Analyze this thumbnail and score its curiosity factor from 0-100.

Thumbnail text: %v
Primary emotion shown: %s

Score based on: information gap, promised reward/transformation, unusual/unexpected element, specific numbers or claims, emotional resonance.

Return ONLY valid JSON: {"curiosity_score": <0-100>, "reasoning": "<2 sentence explanation>"}`, textContent, primaryEmotion)

	reqBody, _ := json.Marshal(map[string]any{
		"contents": []map[string]any{
			{"parts": []map[string]string{{"text": prompt}}},
		},
		"generationConfig": map[string]any{
			"maxOutputTokens": 300,
		},
	})

	endpoint := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", c.model, url.QueryEscape(c.apiKey))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return heuristicCuriosity(textContent, primaryEmotion), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return heuristicCuriosity(textContent, primaryEmotion), nil
	}

	var apiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil || len(apiResp.Candidates) == 0 || len(apiResp.Candidates[0].Content.Parts) == 0 {
		return heuristicCuriosity(textContent, primaryEmotion), nil
	}

	text := stripJSONFence(apiResp.Candidates[0].Content.Parts[0].Text)

	var result CuriosityResult
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return heuristicCuriosity(textContent, primaryEmotion), nil
	}
	return &result, nil
}

var jsonFenceRe = regexp.MustCompile("(?s)```(?:json)?\\s*(.*?)\\s*```")

// stripJSONFence unwraps the markdown code fences Gemini often adds around
// JSON output, even when explicitly asked to return raw JSON.
func stripJSONFence(text string) string {
	text = strings.TrimSpace(text)
	if m := jsonFenceRe.FindStringSubmatch(text); m != nil {
		return strings.TrimSpace(m[1])
	}
	return text
}

var numberRe = regexp.MustCompile(`\d`)

func heuristicCuriosity(textContent []string, primaryEmotion string) *CuriosityResult {
	score := 40
	joined := strings.ToLower(strings.Join(textContent, " "))

	if numberRe.MatchString(joined) {
		score += 15
	}
	for _, hook := range []string{"how", "why", "secret", "never", "before", "after", "vs", "truth", "shocking", "won't"} {
		if strings.Contains(joined, hook) {
			score += 8
			break
		}
	}
	if primaryEmotion == "happy" || primaryEmotion == "surprised" {
		score += 15
	}
	if len(textContent) > 0 {
		score += 10
	}
	if score > 100 {
		score = 100
	}
	return &CuriosityResult{
		Score:     score,
		Reasoning: "Heuristic estimate (no GEMINI_API_KEY configured): based on presence of numbers, curiosity-hook keywords, emotional expression, and text presence.",
	}
}
