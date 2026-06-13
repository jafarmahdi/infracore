package identity

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/infracore/infracore/internal/domain/shared"
	"github.com/jmoiron/sqlx"

	domain "github.com/infracore/infracore/internal/domain/identity"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) domain.UserRepository {
	return &userRepository{db: db}
}

const userSelectCols = `
	id, tenant_id, email, username, password_hash,
	first_name, last_name, phone, avatar_url,
	is_active, is_superuser, email_verified,
	last_login_at, last_login_ip,
	failed_login_attempts, locked_until,
	mfa_enabled, mfa_secret, preferences,
	created_at, updated_at, deleted_at`

func (r *userRepository) Create(ctx context.Context, u *domain.User) error {
	q := `
	INSERT INTO users (
		id, tenant_id, email, username, password_hash,
		first_name, last_name, phone, is_active, is_superuser,
		preferences, created_at, updated_at
	) VALUES (
		:id, :tenant_id, :email, :username, :password_hash,
		:first_name, :last_name, :phone, :is_active, :is_superuser,
		:preferences, :created_at, :updated_at
	)`

	_, err := r.db.NamedExecContext(ctx, q, map[string]any{
		"id":            u.ID,
		"tenant_id":     u.TenantID,
		"email":         u.Email,
		"username":      u.Username,
		"password_hash": u.PasswordHash,
		"first_name":    u.FirstName,
		"last_name":     u.LastName,
		"phone":         u.Phone,
		"is_active":     u.IsActive,
		"is_superuser":  u.IsSuperuser,
		"preferences":   `{}`,
		"created_at":    u.CreatedAt,
		"updated_at":    u.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.User, error) {
	q := fmt.Sprintf(`SELECT %s FROM users WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL`, userSelectCols)
	var row userRow
	if err := r.db.GetContext(ctx, &row, q, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return row.toDomain()
}

func (r *userRepository) GetByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*domain.User, error) {
	q := fmt.Sprintf(`SELECT %s FROM users WHERE tenant_id = $1 AND email = $2 AND deleted_at IS NULL`, userSelectCols)
	var row userRow
	if err := r.db.GetContext(ctx, &row, q, tenantID, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return row.toDomain()
}

func (r *userRepository) GetByUsername(ctx context.Context, tenantID uuid.UUID, username string) (*domain.User, error) {
	q := fmt.Sprintf(`SELECT %s FROM users WHERE tenant_id = $1 AND username = $2 AND deleted_at IS NULL`, userSelectCols)
	var row userRow
	if err := r.db.GetContext(ctx, &row, q, tenantID, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrNotFound
		}
		return nil, fmt.Errorf("get user by username: %w", err)
	}
	return row.toDomain()
}

func (r *userRepository) Update(ctx context.Context, u *domain.User) error {
	q := `
	UPDATE users SET
		first_name = :first_name, last_name = :last_name, phone = :phone,
		avatar_url = :avatar_url, is_active = :is_active,
		email_verified = :email_verified, preferences = :preferences,
		mfa_enabled = :mfa_enabled, mfa_secret = :mfa_secret,
		updated_at = NOW()
	WHERE id = :id AND tenant_id = :tenant_id AND deleted_at IS NULL`

	_, err := r.db.NamedExecContext(ctx, q, map[string]any{
		"id":             u.ID,
		"tenant_id":      u.TenantID,
		"first_name":     u.FirstName,
		"last_name":      u.LastName,
		"phone":          u.Phone,
		"avatar_url":     u.AvatarURL,
		"is_active":      u.IsActive,
		"email_verified": u.EmailVerified,
		"preferences":    `{}`,
		"mfa_enabled":    u.MFAEnabled,
		"mfa_secret":     u.MFASecret,
	})
	return err
}

func (r *userRepository) SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error {
	q := `UPDATE users SET deleted_at = NOW() WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, q, id, tenantID)
	return err
}

func (r *userRepository) List(ctx context.Context, tenantID uuid.UUID, filter domain.UserFilter, page shared.Page) (shared.PageResult[*domain.User], error) {
	q := fmt.Sprintf(`SELECT %s FROM users WHERE tenant_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3`, userSelectCols)
	var rows []userRow
	if err := r.db.SelectContext(ctx, &rows, q, tenantID, page.Size, page.Offset()); err != nil {
		return shared.PageResult[*domain.User]{}, fmt.Errorf("list users: %w", err)
	}

	var total int64
	_ = r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM users WHERE tenant_id = $1 AND deleted_at IS NULL`, tenantID)

	users := make([]*domain.User, 0, len(rows))
	for i := range rows {
		u, err := rows[i].toDomain()
		if err != nil {
			return shared.PageResult[*domain.User]{}, err
		}
		users = append(users, u)
	}
	return shared.NewPageResult(users, total, page), nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, u *domain.User) error {
	q := `
	UPDATE users SET
		last_login_at = $1, last_login_ip = $2,
		failed_login_attempts = 0, locked_until = NULL,
		updated_at = NOW()
	WHERE id = $3`
	_, err := r.db.ExecContext(ctx, q, u.LastLoginAt, u.LastLoginIP.String(), u.ID)
	return err
}
