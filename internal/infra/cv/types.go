package cv

type DominantColor struct {
	Hex        string  `json:"hex"`
	RGB        []int   `json:"rgb"`
	Percentage float64 `json:"percentage"`
	Luminance  float64 `json:"luminance"`
	Saturation float64 `json:"saturation"`
}

type FaceDetail struct {
	BBox            []float64 `json:"bbox"`
	DominantEmotion string    `json:"dominant_emotion"`
	Confidence      float64   `json:"confidence"`
}

type Result struct {
	OCR struct {
		TextDetected     bool     `json:"text_detected"`
		TextStrings      []string `json:"text_strings"`
		TextDensityPct   float64  `json:"text_density_pct"`
		WordCount        int      `json:"word_count"`
		AvgTextHeightPct float64  `json:"avg_text_height_pct"`
	} `json:"ocr"`
	Face struct {
		FaceCount       int          `json:"face_count"`
		HasEyeContact   bool         `json:"has_eye_contact"`
		PrimaryEmotion  string       `json:"primary_emotion"`
		Faces           []FaceDetail `json:"faces"`
	} `json:"face"`
	Colors struct {
		DominantColors  []DominantColor `json:"dominant_colors"`
		ContrastScore   float64         `json:"contrast_score"`
		BrightnessScore float64         `json:"brightness_score"`
		SaturationScore float64         `json:"saturation_score"`
	} `json:"colors"`
	Clutter struct {
		EdgeDensity  float64 `json:"edge_density"`
		ClutterScore float64 `json:"clutter_score"`
		ObjectCount  int     `json:"object_count"`
	} `json:"clutter"`
	VisualComplexity float64 `json:"visual_complexity"`
}
