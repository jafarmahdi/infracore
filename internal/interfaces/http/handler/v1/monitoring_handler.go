package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appmonitoring "github.com/infracore/infracore/internal/application/monitoring"
	"github.com/infracore/infracore/internal/interfaces/http/middleware"
	pkgerrors "github.com/infracore/infracore/pkg/errors"
)

type MonitoringHandler struct {
	svc *appmonitoring.MonitoringService
}

func NewMonitoringHandler(svc *appmonitoring.MonitoringService) *MonitoringHandler {
	return &MonitoringHandler{svc: svc}
}

// ListHosts  GET /api/v1/monitoring/hosts
func (h *MonitoringHandler) ListHosts(c *gin.Context) {
	tenantID := middleware.TenantIDFromCtx(c)
	if tenantID == uuid.Nil {
		pkgerrors.Unauthorized(c, "missing tenant context")
		return
	}

	var req appmonitoring.ListHostsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		pkgerrors.BadRequest(c, "INVALID_QUERY", err.Error())
		return
	}

	resp, err := h.svc.ListHosts(c.Request.Context(), tenantID, req)
	if err != nil {
		pkgerrors.Internal(c)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// CreateHost  POST /api/v1/monitoring/hosts
func (h *MonitoringHandler) CreateHost(c *gin.Context) {
	tenantID := middleware.TenantIDFromCtx(c)
	if tenantID == uuid.Nil {
		pkgerrors.Unauthorized(c, "missing tenant context")
		return
	}

	var req appmonitoring.CreateHostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgerrors.BadRequest(c, "INVALID_BODY", err.Error())
		return
	}

	host, err := h.svc.CreateHost(c.Request.Context(), tenantID, req)
	if err != nil {
		pkgerrors.Internal(c)
		return
	}
	c.JSON(http.StatusCreated, host)
}

// GetHost  GET /api/v1/monitoring/hosts/:id
func (h *MonitoringHandler) GetHost(c *gin.Context) {
	tenantID := middleware.TenantIDFromCtx(c)
	if tenantID == uuid.Nil {
		pkgerrors.Unauthorized(c, "missing tenant context")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkgerrors.BadRequest(c, "INVALID_ID", "invalid host id")
		return
	}

	host, err := h.svc.GetHost(c.Request.Context(), tenantID, id)
	if err != nil {
		pkgerrors.NotFound(c, "host")
		return
	}
	c.JSON(http.StatusOK, host)
}

// DeleteHost  DELETE /api/v1/monitoring/hosts/:id
func (h *MonitoringHandler) DeleteHost(c *gin.Context) {
	tenantID := middleware.TenantIDFromCtx(c)
	if tenantID == uuid.Nil {
		pkgerrors.Unauthorized(c, "missing tenant context")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkgerrors.BadRequest(c, "INVALID_ID", "invalid host id")
		return
	}

	if err := h.svc.DeleteHost(c.Request.Context(), tenantID, id); err != nil {
		pkgerrors.NotFound(c, "host")
		return
	}
	c.Status(http.StatusNoContent)
}

// GetStatusCounts  GET /api/v1/monitoring/hosts/counts
func (h *MonitoringHandler) GetStatusCounts(c *gin.Context) {
	tenantID := middleware.TenantIDFromCtx(c)
	if tenantID == uuid.Nil {
		pkgerrors.Unauthorized(c, "missing tenant context")
		return
	}

	counts, err := h.svc.GetStatusCounts(c.Request.Context(), tenantID)
	if err != nil {
		pkgerrors.Internal(c)
		return
	}
	c.JSON(http.StatusOK, counts)
}
