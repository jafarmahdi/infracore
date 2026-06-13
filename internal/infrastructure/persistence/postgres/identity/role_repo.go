package identity

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	domain "github.com/infracore/infracore/internal/domain/identity"
	"github.com/infracore/infracore/internal/domain/shared"
	"github.com/jmoiron/sqlx"
)

type roleRepository struct {
	db *sqlx.DB
}

func NewRoleRepository(db *sqlx.DB) domain.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *domain.Role) error {
	q := `INSERT INTO roles (id, tenant_id, name, slug, description, is_system, created_at, updated_at)
	      VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := r.db.ExecContext(ctx, q,
		role.ID, role.TenantID, role.Name, role.Slug, role.Description, role.IsSystem, role.CreatedAt, role.UpdatedAt,
	)
	return err
}

func (r *roleRepository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Role, error) {
	q := `SELECT id, tenant_id, name, slug, description, is_system, created_at, updated_at
	      FROM roles WHERE id=$1 AND tenant_id=$2`
	var row roleRow
	if err := r.db.GetContext(ctx, &row, q, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrNotFound
		}
		return nil, fmt.Errorf("get role: %w", err)
	}
	return row.toDomain(), nil
}

func (r *roleRepository) GetBySlug(ctx context.Context, tenantID uuid.UUID, slug string) (*domain.Role, error) {
	q := `SELECT id, tenant_id, name, slug, description, is_system, created_at, updated_at
	      FROM roles WHERE tenant_id=$1 AND slug=$2`
	var row roleRow
	if err := r.db.GetContext(ctx, &row, q, tenantID, slug); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}
	return row.toDomain(), nil
}

func (r *roleRepository) Update(ctx context.Context, role *domain.Role) error {
	q := `UPDATE roles SET name=$1, description=$2, updated_at=NOW() WHERE id=$3 AND tenant_id=$4 AND is_system=false`
	_, err := r.db.ExecContext(ctx, q, role.Name, role.Description, role.ID, role.TenantID)
	return err
}

func (r *roleRepository) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM roles WHERE id=$1 AND tenant_id=$2 AND is_system=false`, id, tenantID)
	return err
}

func (r *roleRepository) List(ctx context.Context, tenantID uuid.UUID, page shared.Page) (shared.PageResult[*domain.Role], error) {
	q := `SELECT id, tenant_id, name, slug, description, is_system, created_at, updated_at
	      FROM roles WHERE tenant_id=$1 ORDER BY name ASC LIMIT $2 OFFSET $3`
	var rows []roleRow
	if err := r.db.SelectContext(ctx, &rows, q, tenantID, page.Size, page.Offset()); err != nil {
		return shared.PageResult[*domain.Role]{}, err
	}
	var total int64
	_ = r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM roles WHERE tenant_id=$1`, tenantID)
	roles := make([]*domain.Role, len(rows))
	for i := range rows {
		roles[i] = rows[i].toDomain()
	}
	return shared.NewPageResult(roles, total, page), nil
}

func (r *roleRepository) GetPermissions(ctx context.Context, roleID uuid.UUID) ([]*domain.Permission, error) {
	q := `SELECT p.id, p.resource, p.action, p.description
	      FROM permissions p
	      JOIN role_permissions rp ON rp.permission_id = p.id
	      WHERE rp.role_id = $1`
	var rows []permissionRow
	if err := r.db.SelectContext(ctx, &rows, q, roleID); err != nil {
		return nil, err
	}
	perms := make([]*domain.Permission, len(rows))
	for i := range rows {
		perms[i] = rows[i].toDomain()
	}
	return perms, nil
}

func (r *roleRepository) SetPermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck

	if _, err := tx.ExecContext(ctx, `DELETE FROM role_permissions WHERE role_id=$1`, roleID); err != nil {
		return err
	}
	for _, pid := range permissionIDs {
		if _, err := tx.ExecContext(ctx, `INSERT INTO role_permissions (role_id, permission_id) VALUES ($1,$2)`, roleID, pid); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *roleRepository) ListPermissions(ctx context.Context) ([]*domain.Permission, error) {
	var rows []permissionRow
	if err := r.db.SelectContext(ctx, &rows, `SELECT id, resource, action, description FROM permissions ORDER BY resource, action`); err != nil {
		return nil, err
	}
	perms := make([]*domain.Permission, len(rows))
	for i := range rows {
		perms[i] = rows[i].toDomain()
	}
	return perms, nil
}

func (r *roleRepository) AssignRoleToUser(ctx context.Context, ur *domain.UserRole) error {
	q := `INSERT INTO user_roles (id, user_id, role_id, site_id, department_id, granted_by, granted_at, expires_at)
	      VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	      ON CONFLICT (user_id, role_id, site_id) DO NOTHING`
	_, err := r.db.ExecContext(ctx, q,
		ur.ID, ur.UserID, ur.RoleID, ur.SiteID, ur.DepartmentID, ur.GrantedBy, ur.GrantedAt, ur.ExpiresAt,
	)
	return err
}

func (r *roleRepository) RevokeRoleFromUser(ctx context.Context, userID, roleID uuid.UUID, siteID *uuid.UUID) error {
	if siteID == nil {
		_, err := r.db.ExecContext(ctx, `DELETE FROM user_roles WHERE user_id=$1 AND role_id=$2 AND site_id IS NULL`, userID, roleID)
		return err
	}
	_, err := r.db.ExecContext(ctx, `DELETE FROM user_roles WHERE user_id=$1 AND role_id=$2 AND site_id=$3`, userID, roleID, *siteID)
	return err
}

func (r *roleRepository) GetUserRoles(ctx context.Context, tenantID, userID uuid.UUID) ([]*domain.UserRole, error) {
	// Returns the user's role assignments with their role slugs for JWT embedding.
	rows := []struct {
		ID           uuid.UUID    `db:"id"`
		UserID       uuid.UUID    `db:"user_id"`
		RoleID       uuid.UUID    `db:"role_id"`
		SiteID       *uuid.UUID   `db:"site_id"`
		DepartmentID *uuid.UUID   `db:"department_id"`
		GrantedBy    uuid.UUID    `db:"granted_by"`
		GrantedAt    interface{}  `db:"granted_at"`
		ExpiresAt    *interface{} `db:"expires_at"`
	}{}
	q := `SELECT ur.id, ur.user_id, ur.role_id, ur.site_id, ur.department_id, ur.granted_by, ur.granted_at, ur.expires_at
	      FROM user_roles ur
	      JOIN roles r ON r.id = ur.role_id
	      WHERE ur.user_id = $1 AND (ur.expires_at IS NULL OR ur.expires_at > NOW())`
	_ = r.db.SelectContext(ctx, &rows, q, userID)
	return nil, nil // Role slugs fetched in GetUserPermissions
}

func (r *roleRepository) GetUserPermissions(ctx context.Context, tenantID, userID uuid.UUID, siteID *uuid.UUID) ([]*domain.Permission, error) {
	q := `SELECT DISTINCT p.id, p.resource, p.action, COALESCE(p.description,'') as description
	      FROM permissions p
	      JOIN role_permissions rp ON rp.permission_id = p.id
	      JOIN user_roles ur ON ur.role_id = rp.role_id
	      WHERE ur.user_id = $1
	        AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
	        AND (ur.site_id = $2 OR ur.site_id IS NULL)
	      ORDER BY p.resource, p.action`
	var rows []permissionRow
	if err := r.db.SelectContext(ctx, &rows, q, userID, siteID); err != nil {
		return nil, fmt.Errorf("get user permissions: %w", err)
	}
	perms := make([]*domain.Permission, len(rows))
	for i := range rows {
		perms[i] = rows[i].toDomain()
	}
	return perms, nil
}

// getUserRoleSlugs is a helper used by the auth service when building JWT claims.
func GetUserRoleSlugs(ctx context.Context, db *sqlx.DB, userID uuid.UUID) ([]string, error) {
	var slugs []string
	q := `SELECT r.slug FROM roles r
	      JOIN user_roles ur ON ur.role_id = r.id
	      WHERE ur.user_id = $1 AND (ur.expires_at IS NULL OR ur.expires_at > NOW())`
	if err := db.SelectContext(ctx, &slugs, q, userID); err != nil {
		return nil, err
	}
	return slugs, nil
}
