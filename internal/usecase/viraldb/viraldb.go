package viraldb

import (
	"context"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/viraldb"
)

type Usecase struct {
	repo viraldb.Repository
}

func NewUsecase(repo viraldb.Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) Search(ctx context.Context, f viraldb.SearchFilter) ([]*viraldb.Thumbnail, error) {
	return u.repo.Search(ctx, f)
}
