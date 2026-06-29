package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "github.com/thumbnailiq/thumbnailiq/pkg/errors"
)

func respondError(c *gin.Context, err error) {
	c.Error(err)
	switch {
	case apperrors.Is(err, apperrors.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	case apperrors.Is(err, apperrors.ErrAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{"error": "already exists"})
	case apperrors.Is(err, apperrors.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	case apperrors.Is(err, apperrors.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	case apperrors.Is(err, apperrors.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
	case apperrors.Is(err, apperrors.ErrQuotaExceeded):
		c.JSON(http.StatusPaymentRequired, gin.H{"error": "analyses quota exceeded for this workspace's plan"})
	case apperrors.Is(err, apperrors.ErrPlanAlreadyActive):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case apperrors.Is(err, apperrors.ErrDowngradeNotAllowed):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
