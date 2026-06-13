package identity

import (
	"database/sql"
	"encoding/json"
	"net"
	"time"

	"github.com/google/uuid"
	domain "github.com/infracore/infracore/internal/domain/identity"
)

// ── DB row models (have sqlx `db` tags, no business logic) ──────

type tenantRow struct {
	ID        uuid.UUID      `db:"id"`
	Name      string         `db:"name"`
	Slug      string         `db:"slug"`
	Plan      string         `db:"plan"`
	MaxUsers  int            `db:"max_users"`
	MaxAssets int            `db:"max_assets"`
	IsActive  bool           `db:"is_active"`
	Settings  []byte         `db:"settings"` // JSONB scanned as []byte
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
	DeletedAt sql.NullTime   `db:"deleted_at"`
}

func (r *tenantRow) toDomain() (*domain.Tenant, error) {
	var settings domain.TenantSettings
	if err := json.Unmarshal(r.Settings, &settings); err != nil {
		return nil, err
	}
	var deletedAt *time.Time
	if r.DeletedAt.Valid {
		t := r.DeletedAt.Time
		deletedAt = &t
	}
	return &domain.Tenant{
		ID:        r.ID,
		Name:      r.Name,
		Slug:      r.Slug,
		Plan:      domain.TenantPlan(r.Plan),
		MaxUsers:  r.MaxUsers,
		MaxAssets: r.MaxAssets,
		IsActive:  r.IsActive,
		Settings:  settings,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		DeletedAt: deletedAt,
	}, nil
}

type userRow struct {
	ID                  uuid.UUID      `db:"id"`
	TenantID            uuid.UUID      `db:"tenant_id"`
	Email               string         `db:"email"`
	Username            string         `db:"username"`
	PasswordHash        string         `db:"password_hash"`
	FirstName           sql.NullString `db:"first_name"`
	LastName            sql.NullString `db:"last_name"`
	Phone               sql.NullString `db:"phone"`
	AvatarURL           sql.NullString `db:"avatar_url"`
	IsActive            bool           `db:"is_active"`
	IsSuperuser         bool           `db:"is_superuser"`
	EmailVerified       bool           `db:"email_verified"`
	LastLoginAt         sql.NullTime   `db:"last_login_at"`
	LastLoginIP         []byte         `db:"last_login_ip"`
	FailedLoginAttempts int            `db:"failed_login_attempts"`
	LockedUntil         sql.NullTime   `db:"locked_until"`
	MFAEnabled          bool           `db:"mfa_enabled"`
	MFASecret           sql.NullString `db:"mfa_secret"`
	Preferences         []byte         `db:"preferences"`
	CreatedAt           time.Time      `db:"created_at"`
	UpdatedAt           time.Time      `db:"updated_at"`
	DeletedAt           sql.NullTime   `db:"deleted_at"`
}

func (r *userRow) toDomain() (*domain.User, error) {
	var prefs domain.UserPreferences
	if err := json.Unmarshal(r.Preferences, &prefs); err != nil {
		return nil, err
	}

	u := &domain.User{
		ID:                  r.ID,
		TenantID:            r.TenantID,
		Email:               r.Email,
		Username:            r.Username,
		PasswordHash:        r.PasswordHash,
		IsActive:            r.IsActive,
		IsSuperuser:         r.IsSuperuser,
		EmailVerified:       r.EmailVerified,
		FailedLoginAttempts: r.FailedLoginAttempts,
		MFAEnabled:          r.MFAEnabled,
		Preferences:         prefs,
		CreatedAt:           r.CreatedAt,
		UpdatedAt:           r.UpdatedAt,
	}

	if r.FirstName.Valid {
		u.FirstName = r.FirstName.String
	}
	if r.LastName.Valid {
		u.LastName = r.LastName.String
	}
	if r.Phone.Valid {
		u.Phone = r.Phone.String
	}
	if r.AvatarURL.Valid {
		u.AvatarURL = r.AvatarURL.String
	}
	if r.MFASecret.Valid {
		u.MFASecret = r.MFASecret.String
	}
	if r.LastLoginAt.Valid {
		t := r.LastLoginAt.Time
		u.LastLoginAt = &t
	}
	if r.LockedUntil.Valid {
		t := r.LockedUntil.Time
		u.LockedUntil = &t
	}
	if len(r.LastLoginIP) > 0 {
		u.LastLoginIP = net.IP(r.LastLoginIP)
	}
	if r.DeletedAt.Valid {
		t := r.DeletedAt.Time
		u.DeletedAt = &t
	}

	return u, nil
}

type permissionRow struct {
	ID          uuid.UUID `db:"id"`
	Resource    string    `db:"resource"`
	Action      string    `db:"action"`
	Description string    `db:"description"`
}

func (r *permissionRow) toDomain() *domain.Permission {
	return &domain.Permission{
		ID:          r.ID,
		Resource:    r.Resource,
		Action:      r.Action,
		Description: r.Description,
	}
}

type roleRow struct {
	ID          uuid.UUID `db:"id"`
	TenantID    uuid.UUID `db:"tenant_id"`
	Name        string    `db:"name"`
	Slug        string    `db:"slug"`
	Description string    `db:"description"`
	IsSystem    bool      `db:"is_system"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (r *roleRow) toDomain() *domain.Role {
	return &domain.Role{
		ID:          r.ID,
		TenantID:    r.TenantID,
		Name:        r.Name,
		Slug:        r.Slug,
		Description: r.Description,
		IsSystem:    r.IsSystem,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

type refreshTokenRow struct {
	ID        uuid.UUID    `db:"id"`
	UserID    uuid.UUID    `db:"user_id"`
	TokenHash string       `db:"token_hash"`
	ExpiresAt time.Time    `db:"expires_at"`
	CreatedAt time.Time    `db:"created_at"`
	IPAddress []byte       `db:"ip_address"`
	UserAgent string       `db:"user_agent"`
	RevokedAt sql.NullTime `db:"revoked_at"`
}

func (r *refreshTokenRow) toDomain() *domain.RefreshToken {
	rt := &domain.RefreshToken{
		ID:        r.ID,
		UserID:    r.UserID,
		TokenHash: r.TokenHash,
		ExpiresAt: r.ExpiresAt,
		CreatedAt: r.CreatedAt,
		UserAgent: r.UserAgent,
	}
	if len(r.IPAddress) > 0 {
		rt.IP = net.IP(r.IPAddress)
	}
	if r.RevokedAt.Valid {
		t := r.RevokedAt.Time
		rt.RevokedAt = &t
	}
	return rt
}
