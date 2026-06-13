package identity

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	domain "github.com/infracore/infracore/internal/domain/identity"
	"github.com/infracore/infracore/internal/domain/shared"
	"github.com/jmoiron/sqlx"
)

type refreshTokenRepository struct {
	db *sqlx.DB
}

func NewRefreshTokenRepository(db *sqlx.DB) domain.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(ctx context.Context, rt *domain.RefreshToken) error {
	q := `
	INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at, ip_address, user_agent)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, q,
		rt.ID, rt.UserID, rt.TokenHash, rt.ExpiresAt, rt.CreatedAt,
		rt.IP.String(), rt.UserAgent,
	)
	return err
}

func (r *refreshTokenRepository) GetByHash(ctx context.Context, hash string) (*domain.RefreshToken, error) {
	q := `SELECT id, user_id, token_hash, expires_at, created_at, ip_address, user_agent, revoked_at
	      FROM refresh_tokens WHERE token_hash = $1`
	var row refreshTokenRow
	if err := r.db.GetContext(ctx, &row, q, hash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrNotFound
		}
		return nil, fmt.Errorf("get refresh token: %w", err)
	}
	return row.toDomain(), nil
}

func (r *refreshTokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE refresh_tokens SET revoked_at=NOW() WHERE id=$1`, id)
	return err
}

func (r *refreshTokenRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE refresh_tokens SET revoked_at=NOW() WHERE user_id=$1 AND revoked_at IS NULL`, userID)
	return err
}

func (r *refreshTokenRepository) DeleteExpired(ctx context.Context) (int64, error) {
	res, err := r.db.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE expires_at < $1`, time.Now())
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
