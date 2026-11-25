package repository

import (
	"context"
	"fmt"
	"strings"

	"LoopForge/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TenantRepositoryInterface 租户仓储接口
type TenantRepositoryInterface interface {
	Create(ctx context.Context, tenant *models.Tenant) error
	GetByID(ctx context.Context, id int64) (*models.Tenant, error)
	Update(ctx context.Context, tenant *models.Tenant) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, req *models.ListTenantRequest) ([]models.Tenant, int64, error)
	UpdateStatus(ctx context.Context, id int64, status int, updateBy string) error
	ExistsByCode(ctx context.Context, code string) (bool, error)
}

// TenantRepository 租户仓储实现
type TenantRepository struct {
	pool *pgxpool.Pool
}

// NewTenantRepository 创建租户仓储
func NewTenantRepository(pool *pgxpool.Pool) *TenantRepository {
	return &TenantRepository{pool: pool}
}

// Create 创建租户
func (r *TenantRepository) Create(ctx context.Context, tenant *models.Tenant) error {
	query := `
		INSERT INTO system_tenant (name, code, domain, status, package_id, expire_time, account_count, create_by, update_by, remark)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, create_time, update_time
	`

	err := r.pool.QueryRow(ctx, query,
		tenant.Name,
		tenant.Code,
		tenant.Domain,
		tenant.Status,
		tenant.PackageID,
		tenant.ExpireTime,
		tenant.AccountCount,
		tenant.CreateBy,
		tenant.UpdateBy,
		tenant.Remark,
	).Scan(&tenant.ID, &tenant.CreateTime, &tenant.UpdateTime)

	if err != nil {
		// 检查唯一约束冲突
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return fmt.Errorf("tenant code already exists: %w", err)
		}
		return fmt.Errorf("create tenant: %w", err)
	}

	return nil
}

// GetByID 根据ID查询租户
func (r *TenantRepository) GetByID(ctx context.Context, id int64) (*models.Tenant, error) {
	query := `
		SELECT id, name, code, domain, status, package_id, expire_time, account_count,
		       create_time, update_time, create_by, update_by, remark
		FROM system_tenant
		WHERE id = $1
	`

	tenant := &models.Tenant{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.Code,
		&tenant.Domain,
		&tenant.Status,
		&tenant.PackageID,
		&tenant.ExpireTime,
		&tenant.AccountCount,
		&tenant.CreateTime,
		&tenant.UpdateTime,
		&tenant.CreateBy,
		&tenant.UpdateBy,
		&tenant.Remark,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("tenant not found: %w", err)
		}
		return nil, fmt.Errorf("get tenant by id: %w", err)
	}

	return tenant, nil
}

// Update 更新租户
func (r *TenantRepository) Update(ctx context.Context, tenant *models.Tenant) error {
	query := `
		UPDATE system_tenant
		SET name = $2, domain = $3, package_id = $4, expire_time = $5,
		    remark = $6, update_by = $7, update_time = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING update_time
	`

	err := r.pool.QueryRow(ctx, query,
		tenant.ID,
		tenant.Name,
		tenant.Domain,
		tenant.PackageID,
		tenant.ExpireTime,
		tenant.Remark,
		tenant.UpdateBy,
	).Scan(&tenant.UpdateTime)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("tenant not found: %w", err)
		}
		return fmt.Errorf("update tenant: %w", err)
	}

	return nil
}

// Delete 删除租户
func (r *TenantRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM system_tenant WHERE id = $1`

	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete tenant: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("tenant not found")
	}

	return nil
}

// List 查询租户列表
func (r *TenantRepository) List(ctx context.Context, req *models.ListTenantRequest) ([]models.Tenant, int64, error) {
	req.SetDefaults()

	// 构建 WHERE 条件
	whereClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Name != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("name ILIKE $%d", argIndex))
		args = append(args, "%"+req.Name+"%")
		argIndex++
	}

	if req.Code != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("code ILIKE $%d", argIndex))
		args = append(args, "%"+req.Code+"%")
		argIndex++
	}

	if req.Status != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *req.Status)
		argIndex++
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM system_tenant %s", whereSQL)
	var total int64
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count tenants: %w", err)
	}

	// 查询列表
	offset := (req.Page - 1) * req.PageSize
	listQuery := fmt.Sprintf(`
		SELECT id, name, code, domain, status, package_id, expire_time, account_count,
		       create_time, update_time, create_by, update_by, remark
		FROM system_tenant
		%s
		ORDER BY create_time DESC
		LIMIT $%d OFFSET $%d
	`, whereSQL, argIndex, argIndex+1)

	args = append(args, req.PageSize, offset)

	rows, err := r.pool.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list tenants: %w", err)
	}
	defer rows.Close()

	tenants := make([]models.Tenant, 0, req.PageSize)
	for rows.Next() {
		var tenant models.Tenant
		err := rows.Scan(
			&tenant.ID,
			&tenant.Name,
			&tenant.Code,
			&tenant.Domain,
			&tenant.Status,
			&tenant.PackageID,
			&tenant.ExpireTime,
			&tenant.AccountCount,
			&tenant.CreateTime,
			&tenant.UpdateTime,
			&tenant.CreateBy,
			&tenant.UpdateBy,
			&tenant.Remark,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan tenant: %w", err)
		}
		tenants = append(tenants, tenant)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate tenants: %w", err)
	}

	return tenants, total, nil
}

// UpdateStatus 更新租户状态
func (r *TenantRepository) UpdateStatus(ctx context.Context, id int64, status int, updateBy string) error {
	query := `
		UPDATE system_tenant
		SET status = $2, update_by = $3, update_time = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING update_time
	`

	var updateTime models.Tenant
	err := r.pool.QueryRow(ctx, query, id, status, updateBy).Scan(&updateTime.UpdateTime)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("tenant not found: %w", err)
		}
		return fmt.Errorf("update tenant status: %w", err)
	}

	return nil
}

// ExistsByCode 检查租户编码是否存在
func (r *TenantRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM system_tenant WHERE code = $1)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, code).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check tenant code exists: %w", err)
	}

	return exists, nil
}
