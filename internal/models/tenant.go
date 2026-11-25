package models

import (
	"time"
)

// Tenant 租户实体
type Tenant struct {
	ID           int64      `json:"id"`                    // 租户ID
	Name         string     `json:"name"`                  // 租户名称
	Code         string     `json:"code"`                  // 租户编码 (唯一)
	Domain       string     `json:"domain,omitempty"`      // 租户域名
	Status       int        `json:"status"`                // 状态: 0-正常, 1-停用
	PackageID    int64      `json:"package_id"`            // 套餐ID
	ExpireTime   *time.Time `json:"expire_time,omitempty"` // 过期时间
	AccountCount int        `json:"account_count"`         // 账号数量
	CreateTime   time.Time  `json:"create_time"`           // 创建时间
	UpdateTime   time.Time  `json:"update_time"`           // 更新时间
	CreateBy     string     `json:"create_by,omitempty"`   // 创建者
	UpdateBy     string     `json:"update_by,omitempty"`   // 更新者
	Remark       string     `json:"remark,omitempty"`      // 备注
}

// CreateTenantRequest 创建租户请求
type CreateTenantRequest struct {
	Name       string     `json:"name" binding:"required,min=2,max=100"`
	Code       string     `json:"code" binding:"required,min=2,max=50"`
	Domain     string     `json:"domain,omitempty" binding:"omitempty,max=255"`
	PackageID  int64      `json:"package_id"`
	ExpireTime *time.Time `json:"expire_time,omitempty"`
	Remark     string     `json:"remark,omitempty" binding:"omitempty,max=500"`
	CreateBy   string     `json:"create_by,omitempty"`
}

// UpdateTenantRequest 更新租户请求
type UpdateTenantRequest struct {
	Name       *string    `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
	Domain     *string    `json:"domain,omitempty" binding:"omitempty,max=255"`
	PackageID  *int64     `json:"package_id,omitempty"`
	ExpireTime *time.Time `json:"expire_time,omitempty"`
	Remark     *string    `json:"remark,omitempty" binding:"omitempty,max=500"`
	UpdateBy   string     `json:"update_by,omitempty"`
}

// ListTenantRequest 查询租户列表请求
type ListTenantRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Name     string `form:"name" binding:"omitempty,max=100"`
	Code     string `form:"code" binding:"omitempty,max=50"`
	Status   *int   `form:"status" binding:"omitempty,oneof=0 1"`
}

// ToggleStatusRequest 切换状态请求
type ToggleStatusRequest struct {
	Status   int    `json:"status" binding:"required,oneof=0 1"`
	UpdateBy string `json:"update_by,omitempty"`
}

// PageResult 分页响应
type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// SetDefaults 设置默认值
func (r *ListTenantRequest) SetDefaults() {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.PageSize <= 0 {
		r.PageSize = 10
	}
	if r.PageSize > 100 {
		r.PageSize = 100
	}
}
