package identity

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/infracore/infracore/internal/domain/shared"
	"github.com/jmoiron/sqlx"

	domain "github.com/infracore/infracore/internal/domain/identity"
)

type tenantRepository struct {
	db *sqlx.DB
}

func NewTenantRepository(db *sqlx.DB) domain.TenantRepository {
	return &tenantRepository{db: db}
}

func (r *tenantRepository) Create(ctx context.Context, t *domain.Tenant) error {
	settingsJSON, err := json.Marshal(t.Settings)
	if err != nil {
		return err
	}
	q := `
	INSERT INTO tenants (id, name, slug, plan, max_users, max_assets, is_active, settings, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err = r.db.ExecContext(ctx, q,
		t.ID, t.Name, t.Slug, string(t.Plan), t.MaxUsers, t.MaxAssets, t.IsActive,
		settingsJSON, t.CreatedAt, t.UpdatedAt,
	)
	return err
}

func (r *tenantRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	q := `SELECT id, name, slug, plan, max_users, max_assets, is_active, settings, created_at, updated_at, deleted_at
	      FROM tenants WHERE id = $1 AND deleted_at IS NULL`
	var row tenantRow
	if err := r.db.GetContext(ctx, &row, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrNotFound
		}
		return nil, fmt.Errorf("get tenant by id: %w", err)
	}
	return row.toDomain()
}

func (r *tenantRepository) GetBySlug(ctx context.Context, slug string) (*domain.Tenant, error) {
	q := `SELECT id, name, slug, plan, max_users, max_assets, is_active, settings, created_at, updated_at, deleted_at
	      FROM tenants WHERE slug = $1 AND deleted_at IS NULL AND is_active = true`
	var row tenantRow
	if err := r.db.GetContext(ctx, &row, q, slug); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrNotFound
		}
		return nil, fmt.Errorf("get tenant by slug: %w", err)
	}
	return row.toDomain()
}

func (r *tenantRepository) Update(ctx context.Context, t *domain.Tenant) error {
	settingsJSON, err := json.Marshal(t.Settings)
	if err != nil {
		return err
	}
	q := `UPDATE tenants SET name=$1, plan=$2, max_users=$3, max_assets=$4, is_active=$5, settings=$6, updated_at=NOW()
	      WHERE id=$7 AND deleted_at IS NULL`
	_, err = r.db.ExecContext(ctx, q, t.Name, string(t.Plan), t.MaxUsers, t.MaxAssets, t.IsActive, settingsJSON, t.ID)
	return err
}

func (r *tenantRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE tenants SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	return err
}

func (r *tenantRepository) List(ctx context.Context, page shared.Page) (shared.PageResult[*domain.Tenant], error) {
	q := `SELECT id, name, slug, plan, max_users, max_assets, is_active, settings, created_at, updated_at, deleted_at
	      FROM tenants WHERE deleted_at IS NULL ORDER BY name ASC LIMIT $1 OFFSET $2`
	var rows []tenantRow
	if err := r.db.SelectContext(ctx, &rows, q, page.Size, page.Offset()); err != nil {
		return shared.PageResult[*domain.Tenant]{}, err
	}
	var total int64
	_ = r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM tenants WHERE deleted_at IS NULL`)

	tenants := make([]*domain.Tenant, 0, len(rows))
	for i := range rows {
		t, err := rows[i].toDomain()
		if err != nil {
			return shared.PageResult[*domain.Tenant]{}, err
		}
		tenants = append(tenants, t)
	}
	return shared.NewPageResult(tenants, total, page), nil
}

func (r *tenantRepository) CountUsers(ctx context.Context, tenantID uuid.UUID) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM users WHERE tenant_id=$1 AND deleted_at IS NULL`, tenantID)
	return count, err
}

func (r *tenantRepository) CountAssets(ctx context.Context, tenantID uuid.UUID) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM assets WHERE tenant_id=$1 AND deleted_at IS NULL`, tenantID)
	return count, err
}

// ensure time package is used
var _ = time.Now
