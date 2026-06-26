package worker

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

const (
	TypeAnalyzeThumbnail = "thumbnail:analyze"
	TypeTrackKeyword     = "tracking:keyword"
	TypeTrackChannel     = "tracking:channel"
)

type AnalyzeThumbnailPayload struct {
	AnalysisID uuid.UUID `json:"analysis_id"`
}

func NewAnalyzeThumbnailTask(analysisID uuid.UUID) (*asynq.Task, error) {
	payload, err := json.Marshal(AnalyzeThumbnailPayload{AnalysisID: analysisID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeAnalyzeThumbnail, payload, asynq.Queue("critical")), nil
}

type TrackingPayload struct {
	TrackingJobID uuid.UUID `json:"tracking_job_id"`
}

func NewTrackKeywordTask(trackingJobID uuid.UUID) (*asynq.Task, error) {
	payload, err := json.Marshal(TrackingPayload{TrackingJobID: trackingJobID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeTrackKeyword, payload, asynq.Queue("default")), nil
}
