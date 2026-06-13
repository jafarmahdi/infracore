package identity

import (
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ─────────────────────────────────────────────────────────────
// Tenant  — aggregate root
// ─────────────────────────────────────────────────────────────

type Tenant struct {
	ID        uuid.UUID
	Name      string
	Slug      string
	Plan      TenantPlan
	MaxUsers  int
	MaxAssets int
	IsActive  bool
	Settings  TenantSettings
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (t *Tenant) IsWithinUserQuota(currentCount int) bool {
	return currentCount < t.MaxUsers
}

func (t *Tenant) IsWithinAssetQuota(currentCount int) bool {
	return currentCount < t.MaxAssets
}

// ─────────────────────────────────────────────────────────────
// Site  — physical branch / location
// ─────────────────────────────────────────────────────────────

type Site struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	Name        string
	Slug        string
	Code        string
	Description string
	Address     Address
	Coordinates *Coordinates
	TimeZone    string
	Contact     Contact
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// ─────────────────────────────────────────────────────────────
// Department
// ─────────────────────────────────────────────────────────────

type Department struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	SiteID      *uuid.UUID
	ParentID    *uuid.UUID
	Name        string
	Description string
	ManagerID   *uuid.UUID
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

func (d *Department) IsRootDepartment() bool {
	return d.ParentID == nil
}

// ─────────────────────────────────────────────────────────────
// User  — entity
// ─────────────────────────────────────────────────────────────

type User struct {
	ID                  uuid.UUID
	TenantID            uuid.UUID
	Email               string
	Username            string
	PasswordHash        string
	FirstName           string
	LastName            string
	Phone               string
	AvatarURL           string
	IsActive            bool
	IsSuperuser         bool
	EmailVerified       bool
	LastLoginAt         *time.Time
	LastLoginIP         net.IP
	FailedLoginAttempts int
	LockedUntil         *time.Time
	MFAEnabled          bool
	MFASecret           string
	Preferences         UserPreferences
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time
}

func (u *User) FullName() string {
	fullName := strings.TrimSpace(strings.Join([]string{u.FirstName, u.LastName}, " "))
	if fullName == "" {
		return u.Username
	}
	return fullName
}

func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

func (u *User) IncrementFailedLogins(maxAttempts int, lockDuration time.Duration) {
	u.FailedLoginAttempts++
	if u.FailedLoginAttempts >= maxAttempts {
		until := time.Now().Add(lockDuration)
		u.LockedUntil = &until
	}
}

func (u *User) RecordLogin(ip net.IP) {
	now := time.Now()
	u.LastLoginAt = &now
	u.LastLoginIP = ip
	u.FailedLoginAttempts = 0
	u.LockedUntil = nil
}

// ─────────────────────────────────────────────────────────────
// Role
// ─────────────────────────────────────────────────────────────

type Role struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	Name        string
	Slug        string
	Description string
	IsSystem    bool // system roles cannot be deleted
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ─────────────────────────────────────────────────────────────
// Permission  — value object
// ─────────────────────────────────────────────────────────────

type Permission struct {
	ID          uuid.UUID
	Resource    string // e.g. "dcim.racks"
	Action      string // create | read | update | delete | list | export
	Description string
}

func (p Permission) Key() string {
	return p.Resource + ":" + p.Action
}

// ─────────────────────────────────────────────────────────────
// UserRole  — assignment of a role to a user, optionally scoped
// ─────────────────────────────────────────────────────────────

type UserRole struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	RoleID       uuid.UUID
	SiteID       *uuid.UUID // nil = tenant-wide
	DepartmentID *uuid.UUID
	GrantedBy    uuid.UUID
	GrantedAt    time.Time
	ExpiresAt    *time.Time
}

func (ur *UserRole) IsExpired() bool {
	if ur.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*ur.ExpiresAt)
}

func (ur *UserRole) IsTenantWide() bool {
	return ur.SiteID == nil
}

// ─────────────────────────────────────────────────────────────
// RefreshToken
// ─────────────────────────────────────────────────────────────

type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string // SHA-256 of the raw token
	ExpiresAt time.Time
	CreatedAt time.Time
	IP        net.IP
	UserAgent string
	RevokedAt *time.Time
}

func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

func (rt *RefreshToken) IsRevoked() bool {
	return rt.RevokedAt != nil
}

func (rt *RefreshToken) IsValid() bool {
	return !rt.IsExpired() && !rt.IsRevoked()
}

// ─────────────────────────────────────────────────────────────
// AuditLog  — immutable change record
// ─────────────────────────────────────────────────────────────

type AuditLog struct {
	ID           uuid.UUID
	TenantID     uuid.UUID
	UserID       *uuid.UUID
	ResourceType string
	ResourceID   *uuid.UUID
	Action       string
	OldValues    map[string]any
	NewValues    map[string]any
	IP           net.IP
	UserAgent    string
	CreatedAt    time.Time
}
