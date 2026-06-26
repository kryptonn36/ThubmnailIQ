package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/viraldb"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/postgres/db"
)

type ViralDBRepo struct {
	q *db.Queries
}

func NewViralDBRepo(pool *pgxpool.Pool) *ViralDBRepo {
	return &ViralDBRepo{q: db.New(pool)}
}

func (r *ViralDBRepo) Search(ctx context.Context, f viraldb.SearchFilter) ([]*viraldb.Thumbnail, error) {
	limit := f.Limit
	if limit <= 0 {
		limit = 50
	}
	var minScore pgtype.Int4
	if f.MinScore != nil {
		minScore = pgtype.Int4{Int32: int32(*f.MinScore), Valid: true}
	}
	var hasFace pgtype.Bool
	if f.HasFace != nil {
		hasFace = pgtype.Bool{Bool: *f.HasFace, Valid: true}
	}
	var niche pgtype.Text
	if f.Niche != nil && *f.Niche != "" {
		niche = pgtype.Text{String: *f.Niche, Valid: true}
	}

	rows, err := r.q.SearchViralThumbnails(ctx, db.SearchViralThumbnailsParams{
		Limit: int32(limit), Offset: int32(f.Offset), Niche: niche, MinScore: minScore, HasFace: hasFace,
	})
	if err != nil {
		return nil, err
	}
	out := make([]*viraldb.Thumbnail, 0, len(rows))
	for _, v := range rows {
		out = append(out, &viraldb.Thumbnail{
			ID: v.ID.String(), VideoID: v.VideoID, ChannelName: v.ChannelName, VideoTitle: v.VideoTitle,
			ThumbnailURL: v.ThumbnailUrl, Niche: textVal(v.Niche), Score: int4Val(v.Score),
			ViewCount: int8Val(v.ViewCount), HasFace: v.HasFace,
		})
	}
	return out, nil
}
