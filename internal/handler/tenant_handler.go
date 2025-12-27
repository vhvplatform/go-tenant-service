package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vhvplatform/go-shared/errors"
	"github.com/vhvplatform/go-shared/logger"
	"github.com/vhvplatform/go-tenant-service/internal/domain"
	"github.com/vhvplatform/go-tenant-service/internal/service"
	"go.uber.org/zap"
)

// TenantHandler handles HTTP requests for tenants
type TenantHandler struct {
	tenantService *service.TenantService
	logger        *logger.Logger
}

// NewTenantHandler creates a new tenant handler
func NewTenantHandler(tenantService *service.TenantService, log *logger.Logger) *TenantHandler {
	return &TenantHandler{
		tenantService: tenantService,
		logger:        log,
	}
}

// CreateTenant handles tenant creation
func (h *TenantHandler) CreateTenant(c *gin.Context) {
	var req domain.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.respondError(c, errors.BadRequest("Invalid request body"))
		return
	}

	tenant, err := h.tenantService.CreateTenant(c.Request.Context(), &req)
	if err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": h.toTenantResponse(tenant)})
}

// GetTenant handles getting a tenant by ID
func (h *TenantHandler) GetTenant(c *gin.Context) {
	tenantID := c.Param("id")

	tenant, err := h.tenantService.GetTenant(c.Request.Context(), tenantID)
	if err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": h.toTenantResponse(tenant)})
}

// ListTenants handles listing tenants
func (h *TenantHandler) ListTenants(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	tenants, total, err := h.tenantService.ListTenants(c.Request.Context(), page, pageSize)
	if err != nil {
		h.respondError(c, err)
		return
	}

	tenantResponses := make([]domain.TenantResponse, len(tenants))
	for i, tenant := range tenants {
		tenantResponses[i] = h.toTenantResponse(tenant)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": domain.ListTenantsResponse{
			Tenants:  tenantResponses,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

// UpdateTenant handles updating a tenant
func (h *TenantHandler) UpdateTenant(c *gin.Context) {
	tenantID := c.Param("id")

	var req domain.UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.respondError(c, errors.BadRequest("Invalid request body"))
		return
	}

	tenant, err := h.tenantService.UpdateTenant(c.Request.Context(), tenantID, &req)
	if err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": h.toTenantResponse(tenant)})
}

// DeleteTenant handles deleting a tenant
func (h *TenantHandler) DeleteTenant(c *gin.Context) {
	tenantID := c.Param("id")

	if err := h.tenantService.DeleteTenant(c.Request.Context(), tenantID); err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenant deleted successfully"})
}

// AddUserToTenant handles adding a user to a tenant
func (h *TenantHandler) AddUserToTenant(c *gin.Context) {
	tenantID := c.Param("id")

	var req domain.AddUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.respondError(c, errors.BadRequest("Invalid request body"))
		return
	}

	if err := h.tenantService.AddUserToTenant(c.Request.Context(), tenantID, req.UserID, req.Role); err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User added to tenant successfully"})
}

// RemoveUserFromTenant handles removing a user from a tenant
func (h *TenantHandler) RemoveUserFromTenant(c *gin.Context) {
	tenantID := c.Param("id")
	userID := c.Param("user_id")

	if err := h.tenantService.RemoveUserFromTenant(c.Request.Context(), tenantID, userID); err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User removed from tenant successfully"})
}

// toTenantResponse converts a tenant domain model to a response
func (h *TenantHandler) toTenantResponse(tenant *domain.Tenant) domain.TenantResponse {
	return domain.TenantResponse{
		ID:               tenant.ID.Hex(),
		Name:             tenant.Name,
		Domain:           tenant.Domain,
		SubscriptionTier: tenant.SubscriptionTier,
		IsActive:         tenant.IsActive,
		Settings:         tenant.Settings,
		CreatedAt:        tenant.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        tenant.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// respondError responds with an error
func (h *TenantHandler) respondError(c *gin.Context, err error) {
	appErr := errors.FromError(err)
	h.logger.Error("Request failed",
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.String("error", appErr.Message),
	)
	c.JSON(appErr.StatusCode, gin.H{"error": appErr})
}

// UpdateTenantConfiguration handles updating tenant configuration
func (h *TenantHandler) UpdateTenantConfiguration(c *gin.Context) {
	tenantID := c.Param("id")

	var req domain.UpdateConfigurationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.respondError(c, errors.BadRequest("Invalid request body"))
		return
	}

	if err := h.tenantService.UpdateTenantConfiguration(c.Request.Context(), tenantID, req.Key, req.Value, req.Type); err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Configuration updated successfully"})
}

// GetTenantConfiguration handles getting tenant configuration
func (h *TenantHandler) GetTenantConfiguration(c *gin.Context) {
	tenantID := c.Param("id")

	config, err := h.tenantService.GetTenantConfiguration(c.Request.Context(), tenantID)
	if err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": config})
}

// GetTenantUsageMetrics handles getting tenant usage metrics
func (h *TenantHandler) GetTenantUsageMetrics(c *gin.Context) {
	tenantID := c.Param("id")
	period := c.DefaultQuery("period", "current")

	metrics, err := h.tenantService.GetTenantUsageMetrics(c.Request.Context(), tenantID, period)
	if err != nil {
		h.respondError(c, err)
		return
	}

	response := domain.UsageMetricsResponse{
		TenantID:      metrics.TenantID,
		APICallCount:  metrics.APICallCount,
		StorageUsed:   metrics.StorageUsed,
		BandwidthUsed: metrics.BandwidthUsed,
		Period:        metrics.Period,
		CreatedAt:     metrics.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// GetTenantUsageHistory handles getting tenant usage history
func (h *TenantHandler) GetTenantUsageHistory(c *gin.Context) {
	tenantID := c.Param("id")
	limit := 30

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	history, err := h.tenantService.GetTenantUsageHistory(c.Request.Context(), tenantID, limit)
	if err != nil {
		h.respondError(c, err)
		return
	}

	var responses []domain.UsageMetricsResponse
	for _, m := range history {
		responses = append(responses, domain.UsageMetricsResponse{
			TenantID:      m.TenantID,
			APICallCount:  m.APICallCount,
			StorageUsed:   m.StorageUsed,
			BandwidthUsed: m.BandwidthUsed,
			Period:        m.Period,
			CreatedAt:     m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}
