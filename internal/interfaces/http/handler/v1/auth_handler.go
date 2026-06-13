package v1

import (
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	appidentity "github.com/infracore/infracore/internal/application/identity"
	"github.com/infracore/infracore/internal/interfaces/http/middleware"
	apierr "github.com/infracore/infracore/pkg/errors"
	"github.com/infracore/infracore/pkg/config"
)

// AuthHandler handles all /auth endpoints.
type AuthHandler struct {
	authSvc *appidentity.AuthService
	cfg     config.AuthConfig
}

func NewAuthHandler(authSvc *appidentity.AuthService, cfg config.AuthConfig) *AuthHandler {
	return &AuthHandler{authSvc: authSvc, cfg: cfg}
}

// Login godoc
// @Summary      Authenticate a user
// @Description  Returns a JWT access token and sets an httpOnly refresh cookie.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      appidentity.LoginRequest   true  "Credentials"
// @Success      200   {object}  appidentity.LoginResponse
// @Failure      400   {object}  errors.ErrorResponse
// @Failure      401   {object}  errors.ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req appidentity.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierr.BadRequest(c, "INVALID_REQUEST", err.Error())
		return
	}

	ip := net.ParseIP(c.ClientIP())
	if ip == nil {
		ip = net.IPv4zero
	}

	resp, rawRefresh, err := h.authSvc.Login(c.Request.Context(), req, ip, c.Request.UserAgent())
	if err != nil {
		switch err.Error() {
		case "[NOT_FOUND] resource not found", "[UNAUTHORIZED] not authorized":
			apierr.Unauthorized(c, "invalid tenant, email, or password")
		case "[FORBIDDEN] access forbidden":
			apierr.Forbidden(c, "account is temporarily locked")
		default:
			apierr.Internal(c)
		}
		return
	}

	h.setRefreshCookie(c, rawRefresh)
	c.JSON(http.StatusOK, resp)
}

// Refresh godoc
// @Summary      Refresh access token
// @Description  Uses the httpOnly refresh cookie to issue a new access token.
// @Tags         auth
// @Produce      json
// @Success      200  {object}  appidentity.RefreshResponse
// @Failure      401  {object}  errors.ErrorResponse
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	rawToken, err := c.Cookie("refresh_token")
	if err != nil || rawToken == "" {
		apierr.Unauthorized(c, "missing refresh token")
		return
	}

	resp, err := h.authSvc.Refresh(c.Request.Context(), rawToken)
	if err != nil {
		h.clearRefreshCookie(c)
		apierr.Unauthorized(c, "invalid or expired refresh token")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Logout godoc
// @Summary      Logout
// @Description  Revokes the refresh token and clears the cookie.
// @Tags         auth
// @Produce      json
// @Success      204
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	rawToken, _ := c.Cookie("refresh_token")
	if rawToken != "" {
		_ = h.authSvc.Logout(c.Request.Context(), rawToken)
	}
	h.clearRefreshCookie(c)
	c.Status(http.StatusNoContent)
}

// Me godoc
// @Summary      Get current user
// @Description  Returns the authenticated user's full profile.
// @Tags         auth
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  appidentity.UserInfo
// @Failure      401  {object}  errors.ErrorResponse
// @Router       /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	tenantID := middleware.TenantIDFromCtx(c)
	userID := middleware.UserIDFromCtx(c)

	info, err := h.authSvc.GetMe(c.Request.Context(), tenantID, userID)
	if err != nil {
		apierr.Internal(c)
		return
	}
	c.JSON(http.StatusOK, info)
}

// ── cookie helpers ───────────────────────────────────────────

func (h *AuthHandler) setRefreshCookie(c *gin.Context, rawToken string) {
	maxAge := int(h.cfg.RefreshTokenTTL.Seconds())
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("refresh_token", rawToken, maxAge, "/api/v1/auth", "", false, true)
}

func (h *AuthHandler) clearRefreshCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("refresh_token", "", -1, "/api/v1/auth", "", false, true)
	_ = time.Now() // suppress unused import
}
