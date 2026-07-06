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
	Keyword  *string
	Niche    *string
	MinScore *int
	HasFace  *bool
	Limit    int
	Offset   int
}

type Repository interface {
	Search(ctx context.Context, f SearchFilter) ([]*Thumbnail, error)
	Upsert(ctx context.Context, t *ThumbnailInput) (*Thumbnail, error)
}

type ThumbnailInput struct {
	VideoID      string
	ChannelID    string
	ChannelName  string
	VideoTitle   string
	ThumbnailURL string
	Niche        string
	Tags         []string
	ViewCount    int64
	Score        int
	HasFace      bool
	CVResults    []byte
}
