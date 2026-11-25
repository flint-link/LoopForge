package service

import (
	"context"
	"fmt"

	"LoopForge/internal/models"
	"LoopForge/internal/repository"
	"go.uber.org/zap"
)

// TenantServiceInterface 租户服务接口
type TenantServiceInterface interface {
	CreateTenant(ctx context.Context, req *models.CreateTenantRequest) (*models.Tenant, error)
	GetTenant(ctx context.Context, id int64) (*models.Tenant, error)
	UpdateTenant(ctx context.Context, id int64, req *models.UpdateTenantRequest) (*models.Tenant, error)
	DeleteTenant(ctx context.Context, id int64) error
	ListTenants(ctx context.Context, req *models.ListTenantRequest) (*models.PageResult, error)
	ToggleTenantStatus(ctx context.Context, id int64, req *models.ToggleStatusRequest) (*models.Tenant, error)
}

// TenantService 租户服务实现
type TenantService struct {
	repo   repository.TenantRepositoryInterface
	logger *zap.Logger
}

// NewTenantService 创建租户服务
func NewTenantService(repo repository.TenantRepositoryInterface, logger *zap.Logger) *TenantService {
	return &TenantService{
		repo:   repo,
		logger: logger,
	}
}

// CreateTenant 创建租户
func (s *TenantService) CreateTenant(ctx context.Context, req *models.CreateTenantRequest) (*models.Tenant, error) {
	s.logger.Info("Creating tenant", zap.String("code", req.Code), zap.String("name", req.Name))

	// 检查编码是否已存在
	exists, err := s.repo.ExistsByCode(ctx, req.Code)
	if err != nil {
		s.logger.Error("Check tenant code exists failed", zap.Error(err))
		return nil, fmt.Errorf("check tenant code: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("tenant code already exists")
	}

	tenant := &models.Tenant{
		Name:         req.Name,
		Code:         req.Code,
		Domain:       req.Domain,
		Status:       0, // 默认正常状态
		PackageID:    req.PackageID,
		ExpireTime:   req.ExpireTime,
		AccountCount: 0,
		CreateBy:     req.CreateBy,
		UpdateBy:     req.CreateBy,
		Remark:       req.Remark,
	}

	if err := s.repo.Create(ctx, tenant); err != nil {
		s.logger.Error("Create tenant failed", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Tenant created successfully", zap.Int64("id", tenant.ID))
	return tenant, nil
}

// GetTenant 获取租户详情
func (s *TenantService) GetTenant(ctx context.Context, id int64) (*models.Tenant, error) {
	s.logger.Debug("Getting tenant", zap.Int64("id", id))

	tenant, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Get tenant failed", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	return tenant, nil
}

// UpdateTenant 更新租户
func (s *TenantService) UpdateTenant(ctx context.Context, id int64, req *models.UpdateTenantRequest) (*models.Tenant, error) {
	s.logger.Info("Updating tenant", zap.Int64("id", id))

	// 先获取现有租户
	tenant, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Get tenant for update failed", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	// 更新字段（只更新非空字段）
	if req.Name != nil {
		tenant.Name = *req.Name
	}
	if req.Domain != nil {
		tenant.Domain = *req.Domain
	}
	if req.PackageID != nil {
		tenant.PackageID = *req.PackageID
	}
	if req.ExpireTime != nil {
		tenant.ExpireTime = req.ExpireTime
	}
	if req.Remark != nil {
		tenant.Remark = *req.Remark
	}
	tenant.UpdateBy = req.UpdateBy

	if err := s.repo.Update(ctx, tenant); err != nil {
		s.logger.Error("Update tenant failed", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	s.logger.Info("Tenant updated successfully", zap.Int64("id", id))
	return tenant, nil
}

// DeleteTenant 删除租户
func (s *TenantService) DeleteTenant(ctx context.Context, id int64) error {
	s.logger.Info("Deleting tenant", zap.Int64("id", id))

	// 先检查租户是否存在
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Get tenant for delete failed", zap.Int64("id", id), zap.Error(err))
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Delete tenant failed", zap.Int64("id", id), zap.Error(err))
		return err
	}

	s.logger.Info("Tenant deleted successfully", zap.Int64("id", id))
	return nil
}

// ListTenants 查询租户列表
func (s *TenantService) ListTenants(ctx context.Context, req *models.ListTenantRequest) (*models.PageResult, error) {
	s.logger.Debug("Listing tenants", zap.Int("page", req.Page), zap.Int("page_size", req.PageSize))

	tenants, total, err := s.repo.List(ctx, req)
	if err != nil {
		s.logger.Error("List tenants failed", zap.Error(err))
		return nil, err
	}

	result := &models.PageResult{
		List:     tenants,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	return result, nil
}

// ToggleTenantStatus 切换租户状态
func (s *TenantService) ToggleTenantStatus(ctx context.Context, id int64, req *models.ToggleStatusRequest) (*models.Tenant, error) {
	s.logger.Info("Toggling tenant status", zap.Int64("id", id), zap.Int("status", req.Status))

	// 先获取现有租户
	tenant, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Get tenant for status toggle failed", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	// 更新状态
	if err := s.repo.UpdateStatus(ctx, id, req.Status, req.UpdateBy); err != nil {
		s.logger.Error("Toggle tenant status failed", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	// 重新获取更新后的租户
	tenant, err = s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Get tenant after status toggle failed", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	s.logger.Info("Tenant status toggled successfully", zap.Int64("id", id), zap.Int("status", req.Status))
	return tenant, nil
}
