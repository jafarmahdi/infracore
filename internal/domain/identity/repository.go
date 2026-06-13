package identity

import (
	"context"

	"github.com/google/uuid"
	"github.com/infracore/infracore/internal/domain/shared"
)

// TenantRepository defines persistence operations for tenants.
type TenantRepository interface {
	Create(ctx context.Context, t *Tenant) error
	GetByID(ctx context.Context, id uuid.UUID) (*Tenant, error)
	GetBySlug(ctx context.Context, slug string) (*Tenant, error)
	Update(ctx context.Context, t *Tenant) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, page shared.Page) (shared.PageResult[*Tenant], error)
	CountUsers(ctx context.Context, tenantID uuid.UUID) (int, error)
	CountAssets(ctx context.Context, tenantID uuid.UUID) (int, error)
}

// SiteRepository defines persistence operations for sites.
type SiteRepository interface {
	Create(ctx context.Context, s *Site) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*Site, error)
	GetBySlug(ctx context.Context, tenantID uuid.UUID, slug string) (*Site, error)
	Update(ctx context.Context, s *Site) error
	SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error
	ListByTenant(ctx context.Context, tenantID uuid.UUID, page shared.Page) (shared.PageResult[*Site], error)
}

// DepartmentRepository defines persistence operations for departments.
type DepartmentRepository interface {
	Create(ctx context.Context, d *Department) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*Department, error)
	Update(ctx context.Context, d *Department) error
	SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error
	ListByTenant(ctx context.Context, tenantID uuid.UUID, page shared.Page) (shared.PageResult[*Department], error)
	ListBySite(ctx context.Context, tenantID, siteID uuid.UUID) ([]*Department, error)
	GetChildren(ctx context.Context, tenantID, parentID uuid.UUID) ([]*Department, error)
}

// UserFilter provides optional filters for user queries.
type UserFilter struct {
	IsActive    *bool
	IsSuperuser *bool
	Search      string // matches email, username, first_name, last_name
}

// UserRepository defines persistence operations for users.
type UserRepository interface {
	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*User, error)
	GetByUsername(ctx context.Context, tenantID uuid.UUID, username string) (*User, error)
	Update(ctx context.Context, u *User) error
	SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error
	List(ctx context.Context, tenantID uuid.UUID, filter UserFilter, page shared.Page) (shared.PageResult[*User], error)
	UpdateLastLogin(ctx context.Context, u *User) error
}

// RoleRepository defines persistence operations for roles and permissions.
type RoleRepository interface {
	Create(ctx context.Context, r *Role) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*Role, error)
	GetBySlug(ctx context.Context, tenantID uuid.UUID, slug string) (*Role, error)
	Update(ctx context.Context, r *Role) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
	List(ctx context.Context, tenantID uuid.UUID, page shared.Page) (shared.PageResult[*Role], error)

	GetPermissions(ctx context.Context, roleID uuid.UUID) ([]*Permission, error)
	SetPermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error
	ListPermissions(ctx context.Context) ([]*Permission, error)

	AssignRoleToUser(ctx context.Context, ur *UserRole) error
	RevokeRoleFromUser(ctx context.Context, userID, roleID uuid.UUID, siteID *uuid.UUID) error
	GetUserRoles(ctx context.Context, tenantID, userID uuid.UUID) ([]*UserRole, error)
	GetUserPermissions(ctx context.Context, tenantID, userID uuid.UUID, siteID *uuid.UUID) ([]*Permission, error)
}

// RefreshTokenRepository manages refresh token lifecycle.
type RefreshTokenRepository interface {
	Create(ctx context.Context, rt *RefreshToken) error
	GetByHash(ctx context.Context, hash string) (*RefreshToken, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeAllForUser(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) (int64, error)
}

// AuditLogRepository writes immutable audit records.
type AuditLogRepository interface {
	Create(ctx context.Context, log *AuditLog) error
	List(ctx context.Context, tenantID uuid.UUID, filter AuditLogFilter, page shared.Page) (shared.PageResult[*AuditLog], error)
}

// AuditLogFilter provides filtering for audit log queries.
type AuditLogFilter struct {
	UserID       *uuid.UUID
	ResourceType string
	ResourceID   *uuid.UUID
	Action       string
}
