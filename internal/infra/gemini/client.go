package gemini

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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

type VideoIdea struct {
	Title       string `json:"title"`
	Hook        string `json:"hook"`
	Format      string `json:"format"`
	CTRPotential string `json:"ctr_potential"`
}

type VideoIdeasResult struct {
	Ideas []VideoIdea `json:"ideas"`
}

// GenerateVideoIdeas suggests fresh video ideas for a keyword by studying
// what titles and formats are already working for the top-ranking competitors.
// Falls back to a rule-based template generator when Gemini isn't available.
func (c *Client) GenerateVideoIdeas(ctx context.Context, keyword string, competitorTitles []string) (*VideoIdeasResult, error) {
	if c.apiKey == "" || len(competitorTitles) == 0 {
		return heuristicIdeas(keyword, competitorTitles), nil
	}

	titlesStr := ""
	for i, t := range competitorTitles {
		if i >= 15 {
			break
		}
		titlesStr += fmt.Sprintf("- %s\n", t)
	}

	prompt := fmt.Sprintf(`You are a top YouTube content strategist. A creator wants to make a video about %q.

Here are the top-performing competitor titles in this niche:
%s

Generate exactly 6 fresh, original video ideas that would compete well against these. Each idea must:
- Have a distinct format (e.g. listicle, how-to, story, challenge, comparison, case study)
- Use a proven click-trigger (curiosity gap, number, emotion, controversy, FOMO, transformation)
- Be genuinely different from the competitors above

Return ONLY valid JSON:
{"ideas": [
  {"title": "<catchy title>", "hook": "<1-sentence hook explaining why viewers click>", "format": "<Listicle|How-To|Story|Challenge|Comparison|Case Study|Reaction>", "ctr_potential": "<Low|Medium|High|Very High>"},
  ...
]}`, keyword, titlesStr)

	reqBody, _ := json.Marshal(map[string]any{
		"contents": []map[string]any{
			{"parts": []map[string]string{{"text": prompt}}},
		},
		"generationConfig": map[string]any{"maxOutputTokens": 1200},
	})

	endpoint := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", c.model, url.QueryEscape(c.apiKey))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return heuristicIdeas(keyword, competitorTitles), nil
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return heuristicIdeas(keyword, competitorTitles), nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return heuristicIdeas(keyword, competitorTitles), nil
	}

	var apiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct{ Text string `json:"text"` } `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil || len(apiResp.Candidates) == 0 || len(apiResp.Candidates[0].Content.Parts) == 0 {
		return heuristicIdeas(keyword, competitorTitles), nil
	}
	text := stripJSONFence(apiResp.Candidates[0].Content.Parts[0].Text)
	var result VideoIdeasResult
	if err := json.Unmarshal([]byte(text), &result); err != nil || len(result.Ideas) == 0 {
		return heuristicIdeas(keyword, competitorTitles), nil
	}
	return &result, nil
}

func heuristicIdeas(keyword string, _ []string) *VideoIdeasResult {
	kw := keyword
	templates := []VideoIdea{
		{Title: fmt.Sprintf("I Tried %s For 30 Days — Here's What Happened", strings.Title(kw)), Hook: "Transformation story triggers curiosity about the result", Format: "Story", CTRPotential: "High"},
		{Title: fmt.Sprintf("7 %s Mistakes Everyone Makes (And How To Fix Them)", strings.Title(kw)), Hook: "Listicle + fear of being wrong makes viewers click to self-check", Format: "Listicle", CTRPotential: "Very High"},
		{Title: fmt.Sprintf("The Truth About %s Nobody Talks About", strings.Title(kw)), Hook: "Curiosity gap — implies hidden or suppressed information", Format: "Story", CTRPotential: "High"},
		{Title: fmt.Sprintf("How To %s Even If You're a Complete Beginner", strings.Title(kw)), Hook: "Low barrier framing attracts the largest possible audience", Format: "How-To", CTRPotential: "Very High"},
		{Title: fmt.Sprintf("%s vs %s: I Tested Both For 60 Days", strings.Title(kw), strings.Title(kw)+" Pro"), Hook: "Comparison + personal test beats generic reviews for CTR", Format: "Comparison", CTRPotential: "High"},
		{Title: fmt.Sprintf("Why I Quit %s (And What I Do Instead)", strings.Title(kw)), Hook: "Contrarian hook builds curiosity and emotional investment", Format: "Story", CTRPotential: "Very High"},
	}
	return &VideoIdeasResult{Ideas: templates}
}

type RelevanceResult struct {
	Score     int    `json:"relevance_score"`
	Reasoning string `json:"reasoning"`
}

// ScoreRelevance judges whether a thumbnail's actual visual content has
// anything to do with the keyword it was uploaded for — a check the CV
// pipeline never performs on its own, since every other sub-score only
// looks at the image's intrinsic properties (contrast, faces, clutter...)
// compared against competitors, never against the search topic. A
// visually polished but completely off-topic thumbnail would otherwise
// rank well purely on those generic heuristics.
//
// Without a working Gemini key, this falls back to a much weaker
// OCR-text/keyword overlap heuristic — it can't see the image, only
// whatever text the CV pipeline already extracted from it.
func (c *Client) ScoreRelevance(ctx context.Context, imageURL, keyword string, ocrText []string) (*RelevanceResult, error) {
	if c.apiKey == "" {
		return heuristicRelevance(ocrText, keyword), nil
	}

	imgBytes, mimeType, err := fetchImage(ctx, imageURL)
	if err != nil {
		return heuristicRelevance(ocrText, keyword), nil
	}

	prompt := fmt.Sprintf(`You are judging whether a YouTube thumbnail image is topically relevant to a search keyword.

Search keyword: %q

Look at the actual image content (objects, people, text, scene) and judge how relevant it is to that keyword. A thumbnail that has nothing to do with the keyword should score low even if it's visually polished — this is purely a topical-relevance check, not a quality check.

Return ONLY valid JSON: {"relevance_score": <0-100>, "reasoning": "<1-2 sentence explanation>"}`, keyword)

	reqBody, _ := json.Marshal(map[string]any{
		"contents": []map[string]any{
			{"parts": []map[string]any{
				{"text": prompt},
				{"inlineData": map[string]string{"mimeType": mimeType, "data": base64.StdEncoding.EncodeToString(imgBytes)}},
			}},
		},
		"generationConfig": map[string]any{
			"maxOutputTokens": 200,
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
		return heuristicRelevance(ocrText, keyword), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return heuristicRelevance(ocrText, keyword), nil
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
		return heuristicRelevance(ocrText, keyword), nil
	}

	text := stripJSONFence(apiResp.Candidates[0].Content.Parts[0].Text)

	var result RelevanceResult
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return heuristicRelevance(ocrText, keyword), nil
	}
	return &result, nil
}

func fetchImage(ctx context.Context, imageURL string) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		return nil, "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("fetch image: status %d", resp.StatusCode)
	}
	mimeType := resp.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "image/jpeg"
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	return body, mimeType, nil
}

// heuristicRelevance can't see the image at all — it only knows whatever
// text the CV pipeline's OCR step already extracted, so it's a much weaker
// signal than the real vision-based check and is clearly labeled as such.
func heuristicRelevance(ocrText []string, keyword string) *RelevanceResult {
	words := strings.Fields(strings.ToLower(keyword))
	if len(words) == 0 {
		return &RelevanceResult{Score: 100, Reasoning: "No keyword to check relevance against."}
	}

	joined := strings.ToLower(strings.Join(ocrText, " "))
	matches := 0
	for _, w := range words {
		if len(w) > 2 && strings.Contains(joined, w) {
			matches++
		}
	}

	score := 30 + int(70*float64(matches)/float64(len(words)))
	return &RelevanceResult{
		Score: score,
		Reasoning: "Heuristic estimate (no working GEMINI_API_KEY): based only on whether the keyword's words appear in the " +
			"thumbnail's OCR text, not the actual visual content. Configure a real Gemini key for an accurate check.",
	}
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
