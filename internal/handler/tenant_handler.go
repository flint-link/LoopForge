package handler

import (
	"net/http"
	"strconv"

	"LoopForge/internal/models"
	"LoopForge/internal/response"
	"LoopForge/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TenantHandler 租户处理器
type TenantHandler struct {
	service service.TenantServiceInterface
	logger  *zap.Logger
}

// NewTenantHandler 创建租户处理器
func NewTenantHandler(service service.TenantServiceInterface, logger *zap.Logger) *TenantHandler {
	return &TenantHandler{
		service: service,
		logger:  logger,
	}
}

// Create 创建租户
// @Summary 创建租户
// @Description 创建新租户
// @Tags 租户管理
// @Accept json
// @Produce json
// @Param request body models.CreateTenantRequest true "创建租户请求"
// @Success 200 {object} response.Response{data=models.Tenant}
// @Router /api/v1/tenants [post]
func (h *TenantHandler) Create(c *gin.Context) {
	var req models.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	tenant, err := h.service.CreateTenant(c.Request.Context(), &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, tenant)
}

// List 查询租户列表
// @Summary 查询租户列表
// @Description 分页查询租户列表，支持按名称、编码、状态筛选
// @Tags 租户管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param name query string false "租户名称"
// @Param code query string false "租户编码"
// @Param status query int false "状态: 0-正常, 1-停用"
// @Success 200 {object} response.Response{data=models.PageResult}
// @Router /api/v1/tenants [get]
func (h *TenantHandler) List(c *gin.Context) {
	var req models.ListTenantRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	result, err := h.service.ListTenants(c.Request.Context(), &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, result)
}

// Get 查询租户详情
// @Summary 查询租户详情
// @Description 根据ID查询租户详情
// @Tags 租户管理
// @Accept json
// @Produce json
// @Param id path int true "租户ID"
// @Success 200 {object} response.Response{data=models.Tenant}
// @Router /api/v1/tenants/{id} [get]
func (h *TenantHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeBadRequest, "Invalid tenant ID")
		return
	}

	tenant, err := h.service.GetTenant(c.Request.Context(), id)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, tenant)
}

// Update 更新租户
// @Summary 更新租户
// @Description 更新租户信息（部分字段更新）
// @Tags 租户管理
// @Accept json
// @Produce json
// @Param id path int true "租户ID"
// @Param request body models.UpdateTenantRequest true "更新租户请求"
// @Success 200 {object} response.Response{data=models.Tenant}
// @Router /api/v1/tenants/{id} [put]
func (h *TenantHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeBadRequest, "Invalid tenant ID")
		return
	}

	var req models.UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	tenant, err := h.service.UpdateTenant(c.Request.Context(), id, &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, tenant)
}

// Delete 删除租户
// @Summary 删除租户
// @Description 删除租户（硬删除）
// @Tags 租户管理
// @Accept json
// @Produce json
// @Param id path int true "租户ID"
// @Success 200 {object} response.Response
// @Router /api/v1/tenants/{id} [delete]
func (h *TenantHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeBadRequest, "Invalid tenant ID")
		return
	}

	err = h.service.DeleteTenant(c.Request.Context(), id)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Tenant deleted successfully"})
}

// ToggleStatus 切换租户状态
// @Summary 切换租户状态
// @Description 启用/禁用租户
// @Tags 租户管理
// @Accept json
// @Produce json
// @Param id path int true "租户ID"
// @Param request body models.ToggleStatusRequest true "切换状态请求"
// @Success 200 {object} response.Response{data=models.Tenant}
// @Router /api/v1/tenants/{id}/status [patch]
func (h *TenantHandler) ToggleStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeBadRequest, "Invalid tenant ID")
		return
	}

	var req models.ToggleStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	tenant, err := h.service.ToggleTenantStatus(c.Request.Context(), id, &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, tenant)
}

// RegisterRoutes 注册路由
func (h *TenantHandler) RegisterRoutes(router *gin.RouterGroup) {
	tenants := router.Group("/tenants")
	{
		tenants.POST("", h.Create)
		tenants.GET("", h.List)
		tenants.GET("/:id", h.Get)
		tenants.PUT("/:id", h.Update)
		tenants.DELETE("/:id", h.Delete)
		tenants.PATCH("/:id/status", h.ToggleStatus)
	}
}

// HealthCheck 健康检查
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data":    "healthy",
	})
}
