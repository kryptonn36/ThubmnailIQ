package scoring

// CompetitorAvg holds the mean of each numeric CV metric across the
// competitor set, used as the niche benchmark for relative scoring.
type CompetitorAvg struct {
	Score          float64 `json:"score"`
	FaceCount      float64 `json:"face_count"`
	TextDensityPct float64 `json:"text_density_pct"`
	ClutterScore   float64 `json:"clutter_score"`
	ContrastScore  float64 `json:"contrast_score"`
	SaturationScore float64 `json:"saturation_score"`
}

type SubScores struct {
	Visibility int
	Contrast   int
	Attention  int
	Mobile     int
	Branding   int
	Curiosity  int
}

type Suggestion struct {
	Type             string  `json:"type"`
	ImpactLevel      string  `json:"impact_level"`
	Headline         string  `json:"headline"`
	Explanation      string  `json:"explanation"`
	EstimatedCTRMin  float64 `json:"estimated_ctr_min"`
	EstimatedCTRMax  float64 `json:"estimated_ctr_max"`
}

func clamp(v float64, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
