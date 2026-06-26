package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/thumbnailiq/thumbnailiq/internal/middleware"
	workspaceuc "github.com/thumbnailiq/thumbnailiq/internal/usecase/workspace"
	apperrors "github.com/thumbnailiq/thumbnailiq/pkg/errors"
)

// resolveWorkspaceID lets clients omit workspace_id and fall back to the
// caller's first workspace. Every user gets exactly one auto-created
// workspace at registration in this build, so this keeps single-workspace
// clients (like the bundled web app) simple while still letting
// multi-workspace API consumers pass workspace_id explicitly.
func resolveWorkspaceID(c *gin.Context, raw string, workspaces *workspaceuc.Usecase) (uuid.UUID, error) {
	if raw != "" {
		return uuid.Parse(raw)
	}
	list, err := workspaces.List(c.Request.Context(), middleware.UserID(c))
	if err != nil {
		return uuid.UUID{}, err
	}
	if len(list) == 0 {
		return uuid.UUID{}, apperrors.ErrNotFound
	}
	return list[0].ID, nil
}
