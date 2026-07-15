package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/thumbnailiq/thumbnailiq/internal/middleware"
	workspaceuc "github.com/thumbnailiq/thumbnailiq/internal/usecase/workspace"
	apperrors "github.com/thumbnailiq/thumbnailiq/pkg/errors"
)

// resolveWorkspaceID resolves which workspace a request operates on, most
// specific source first: an explicit workspace_id in the query/body, then the
// X-Workspace-ID header (set by the web app to its currently selected
// workspace, so every endpoint follows the switcher without each call site
// threading an id), and finally the caller's first workspace as a fallback
// for single-workspace clients.
func resolveWorkspaceID(c *gin.Context, raw string, workspaces *workspaceuc.Usecase) (uuid.UUID, error) {
	if raw == "" {
		raw = c.GetHeader("X-Workspace-ID")
	}
	if raw != "" {
		id, err := uuid.Parse(raw)
		if err != nil {
			return uuid.UUID{}, err
		}
		// A client-supplied workspace ID is never trusted: confirm the caller
		// is actually a member before any handler acts on it. Without this,
		// knowing another workspace's UUID would be enough to operate on it.
		if err := workspaces.EnsureMember(c.Request.Context(), middleware.UserID(c), id); err != nil {
			return uuid.UUID{}, err
		}
		return id, nil
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
