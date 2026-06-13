package identity

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "github.com/infracore/infracore/internal/domain/identity"
	"github.com/infracore/infracore/internal/domain/shared"
	pgidentity "github.com/infracore/infracore/internal/infrastructure/persistence/postgres/identity"
	"github.com/infracore/infracore/pkg/config"
	"github.com/infracore/infracore/pkg/crypto"
	"github.com/jmoiron/sqlx"
)

// ── DTOs ────────────────────────────────────────────────────────

type LoginRequest struct {
	TenantSlug string `json:"tenant_slug" binding:"required"`
	Email      string `json:"email"       binding:"required,email"`
	Password   string `json:"password"    binding:"required,min=6"`
}

type LoginResponse struct {
	AccessToken string   `json:"access_token"`
	TokenType   string   `json:"token_type"`
	ExpiresIn   int      `json:"expires_in"` // seconds
	User        UserInfo `json:"user"`
}

type UserInfo struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	TenantSlug  string    `json:"tenant_slug"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	AvatarURL   string    `json:"avatar_url"`
	IsSuperuser bool      `json:"is_superuser"`
	Roles       []string  `json:"roles"`
	Permissions []string  `json:"permissions"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// ── Service ─────────────────────────────────────────────────────

type AuthService struct {
	tenantRepo domain.TenantRepository
	userRepo   domain.UserRepository
	roleRepo   domain.RoleRepository
	tokenRepo  domain.RefreshTokenRepository
	jwtMgr     *crypto.JWTManager
	cfg        config.AuthConfig
	log        *zap.Logger
	db         *sqlx.DB
}

func NewAuthService(
	tenantRepo domain.TenantRepository,
	userRepo domain.UserRepository,
	roleRepo domain.RoleRepository,
	tokenRepo domain.RefreshTokenRepository,
	jwtMgr *crypto.JWTManager,
	cfg config.AuthConfig,
	log *zap.Logger,
	db *sqlx.DB,
) *AuthService {
	return &AuthService{
		tenantRepo: tenantRepo,
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		tokenRepo:  tokenRepo,
		jwtMgr:     jwtMgr,
		cfg:        cfg,
		log:        log,
		db:         db,
	}
}

// Login authenticates a user and returns tokens.
func (s *AuthService) Login(ctx context.Context, req LoginRequest, ip net.IP, userAgent string) (*LoginResponse, string, error) {
	// 1. Resolve tenant
	tenant, err := s.tenantRepo.GetBySlug(ctx, req.TenantSlug)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return nil, "", shared.ErrNotFound
		}
		return nil, "", fmt.Errorf("resolve tenant: %w", err)
	}

	// 2. Find user
	user, err := s.userRepo.GetByEmail(ctx, tenant.ID, req.Email)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return nil, "", shared.ErrUnauthorized
		}
		return nil, "", fmt.Errorf("find user: %w", err)
	}

	// 3. Check account state
	if !user.IsActive {
		return nil, "", shared.ErrUnauthorized
	}
	if user.IsLocked() {
		return nil, "", shared.ErrForbidden
	}

	// 4. Verify password
	if !crypto.CheckPassword(req.Password, user.PasswordHash) {
		user.IncrementFailedLogins(s.cfg.MaxFailedLogins, s.cfg.LockoutDuration)
		_ = s.userRepo.Update(ctx, user)
		return nil, "", shared.ErrUnauthorized
	}

	// 5. Fetch roles + permissions for JWT
	roleSlugs, err := pgidentity.GetUserRoleSlugs(ctx, s.db, user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("fetch roles: %w", err)
	}
	permissions, err := s.roleRepo.GetUserPermissions(ctx, tenant.ID, user.ID, nil)
	if err != nil {
		return nil, "", fmt.Errorf("fetch permissions: %w", err)
	}
	permStrings := make([]string, len(permissions))
	for i, p := range permissions {
		permStrings[i] = p.Key()
	}

	// 6. Generate access token
	accessToken, err := s.jwtMgr.GenerateAccessToken(crypto.AccessTokenClaims{
		UserID:      user.ID,
		TenantID:    tenant.ID,
		TenantSlug:  tenant.Slug,
		Email:       user.Email,
		Username:    user.Username,
		IsSuperuser: user.IsSuperuser,
		Roles:       roleSlugs,
		Permissions: permStrings,
	})
	if err != nil {
		return nil, "", fmt.Errorf("generate access token: %w", err)
	}

	// 7. Generate refresh token
	rawRefresh, err := crypto.GenerateSecureToken(32)
	if err != nil {
		return nil, "", err
	}
	refreshHash := crypto.SHA256Hex(rawRefresh)
	now := time.Now()
	rt := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: refreshHash,
		ExpiresAt: now.Add(s.jwtMgr.RefreshTTL()),
		CreatedAt: now,
		IP:        ip,
		UserAgent: userAgent,
	}
	if err := s.tokenRepo.Create(ctx, rt); err != nil {
		return nil, "", fmt.Errorf("save refresh token: %w", err)
	}

	// 8. Record login
	user.RecordLogin(ip)
	_ = s.userRepo.UpdateLastLogin(ctx, user)

	s.log.Info("user logged in",
		zap.String("user_id", user.ID.String()),
		zap.String("tenant", tenant.Slug),
		zap.String("ip", ip.String()),
	)

	resp := &LoginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(s.cfg.AccessTokenTTL.Seconds()),
		User: UserInfo{
			ID:          user.ID,
			TenantID:    tenant.ID,
			TenantSlug:  tenant.Slug,
			Email:       user.Email,
			Username:    user.Username,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			AvatarURL:   user.AvatarURL,
			IsSuperuser: user.IsSuperuser,
			Roles:       roleSlugs,
			Permissions: permStrings,
		},
	}
	return resp, rawRefresh, nil
}

// Refresh issues a new access token given a valid refresh token.
func (s *AuthService) Refresh(ctx context.Context, rawRefreshToken string) (*RefreshResponse, error) {
	hash := crypto.SHA256Hex(rawRefreshToken)
	rt, err := s.tokenRepo.GetByHash(ctx, hash)
	if err != nil {
		return nil, shared.ErrUnauthorized
	}
	if !rt.IsValid() {
		return nil, shared.ErrUnauthorized
	}

	user, err := s.userRepo.GetByID(ctx, uuid.Nil, rt.UserID) // tenant_id not enforced here
	if err != nil {
		return nil, shared.ErrUnauthorized
	}

	tenant, err := s.tenantRepo.GetByID(ctx, user.TenantID)
	if err != nil {
		return nil, shared.ErrUnauthorized
	}

	roleSlugs, _ := pgidentity.GetUserRoleSlugs(ctx, s.db, user.ID)
	permissions, _ := s.roleRepo.GetUserPermissions(ctx, tenant.ID, user.ID, nil)
	permStrings := make([]string, len(permissions))
	for i, p := range permissions {
		permStrings[i] = p.Key()
	}

	accessToken, err := s.jwtMgr.GenerateAccessToken(crypto.AccessTokenClaims{
		UserID:      user.ID,
		TenantID:    tenant.ID,
		TenantSlug:  tenant.Slug,
		Email:       user.Email,
		Username:    user.Username,
		IsSuperuser: user.IsSuperuser,
		Roles:       roleSlugs,
		Permissions: permStrings,
	})
	if err != nil {
		return nil, err
	}

	return &RefreshResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(s.cfg.AccessTokenTTL.Seconds()),
	}, nil
}

// Logout revokes the provided refresh token.
func (s *AuthService) Logout(ctx context.Context, rawRefreshToken string) error {
	hash := crypto.SHA256Hex(rawRefreshToken)
	rt, err := s.tokenRepo.GetByHash(ctx, hash)
	if err != nil {
		return nil // already gone — treat as success
	}
	return s.tokenRepo.Revoke(ctx, rt.ID)
}

// GetMe returns the current user's full profile from the DB.
func (s *AuthService) GetMe(ctx context.Context, tenantID, userID uuid.UUID) (*UserInfo, error) {
	user, err := s.userRepo.GetByID(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	roleSlugs, _ := pgidentity.GetUserRoleSlugs(ctx, s.db, userID)
	permissions, _ := s.roleRepo.GetUserPermissions(ctx, tenantID, userID, nil)
	permStrings := make([]string, len(permissions))
	for i, p := range permissions {
		permStrings[i] = p.Key()
	}
	return &UserInfo{
		ID:          user.ID,
		TenantID:    tenant.ID,
		TenantSlug:  tenant.Slug,
		Email:       user.Email,
		Username:    user.Username,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		AvatarURL:   user.AvatarURL,
		IsSuperuser: user.IsSuperuser,
		Roles:       roleSlugs,
		Permissions: permStrings,
	}, nil
}
