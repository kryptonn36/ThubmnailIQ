package viraldb

import "context"

type Thumbnail struct {
	ID           string `json:"id"`
	VideoID      string `json:"video_id"`
	ChannelName  string `json:"channel_name"`
	VideoTitle   string `json:"video_title"`
	ThumbnailURL string `json:"thumbnail_url"`
	Niche        string `json:"niche"`
	Score        int    `json:"score"`
	ViewCount    int64  `json:"view_count"`
	HasFace      bool   `json:"has_face"`
}

type SearchFilter struct {
	Niche    *string
	MinScore *int
	HasFace  *bool
	Limit    int
	Offset   int
}

type Repository interface {
	Search(ctx context.Context, f SearchFilter) ([]*Thumbnail, error)
}
