package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	apierr "github.com/infracore/infracore/pkg/errors"
	"github.com/infracore/infracore/pkg/crypto"
	"github.com/google/uuid"
)

// ContextKey constants for values stored in Gin's context.
const (
	CtxUserID      = "user_id"
	CtxTenantID    = "tenant_id"
	CtxTenantSlug  = "tenant_slug"
	CtxEmail       = "email"
	CtxUsername    = "username"
	CtxIsSuperuser = "is_superuser"
	CtxRoles       = "roles"
	CtxPermissions = "permissions"
)

// Auth returns middleware that validates the Bearer JWT in Authorization header.
func Auth(jwtMgr *crypto.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			apierr.AbortUnauthorized(c, "missing Authorization header")
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			apierr.AbortUnauthorized(c, "invalid Authorization header format")
			return
		}

		claims, err := jwtMgr.ParseAccessToken(parts[1])
		if err != nil {
			apierr.AbortUnauthorized(c, "invalid or expired token")
			return
		}

		// Inject caller identity into request context
		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxTenantID, claims.TenantID)
		c.Set(CtxTenantSlug, claims.TenantSlug)
		c.Set(CtxEmail, claims.Email)
		c.Set(CtxUsername, claims.Username)
		c.Set(CtxIsSuperuser, claims.IsSuperuser)
		c.Set(CtxRoles, claims.Roles)
		c.Set(CtxPermissions, claims.Permissions)

		c.Next()
	}
}

// RequirePermission returns middleware that checks the caller has a specific permission.
// Superusers bypass all permission checks.
func RequirePermission(resource, action string) gin.HandlerFunc {
	required := resource + ":" + action
	return func(c *gin.Context) {
		if isSuperuser, _ := c.Get(CtxIsSuperuser); isSuperuser == true {
			c.Next()
			return
		}
		permsRaw, exists := c.Get(CtxPermissions)
		if !exists {
			apierr.AbortForbidden(c, "insufficient permissions")
			return
		}
		perms, ok := permsRaw.([]string)
		if !ok {
			apierr.AbortForbidden(c, "insufficient permissions")
			return
		}
		for _, p := range perms {
			if p == required {
				c.Next()
				return
			}
		}
		apierr.AbortForbidden(c, "you lack the '"+required+"' permission")
	}
}

// RequireRole returns middleware that checks the caller holds at least one of the given role slugs.
func RequireRole(slugs ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(slugs))
	for _, s := range slugs {
		allowed[s] = struct{}{}
	}
	return func(c *gin.Context) {
		if isSuperuser, _ := c.Get(CtxIsSuperuser); isSuperuser == true {
			c.Next()
			return
		}
		rolesRaw, exists := c.Get(CtxRoles)
		if !exists {
			apierr.AbortForbidden(c, "insufficient role")
			return
		}
		roles, ok := rolesRaw.([]string)
		if !ok {
			apierr.AbortForbidden(c, "insufficient role")
			return
		}
		for _, r := range roles {
			if _, ok := allowed[r]; ok {
				c.Next()
				return
			}
		}
		apierr.AbortForbidden(c, "insufficient role")
	}
}

// ── Context extraction helpers used by handlers ───────────────

func UserIDFromCtx(c *gin.Context) uuid.UUID {
	if v, ok := c.Get(CtxUserID); ok {
		if id, ok := v.(uuid.UUID); ok {
			return id
		}
	}
	return uuid.Nil
}

func TenantIDFromCtx(c *gin.Context) uuid.UUID {
	if v, ok := c.Get(CtxTenantID); ok {
		if id, ok := v.(uuid.UUID); ok {
			return id
		}
	}
	return uuid.Nil
}

func TenantSlugFromCtx(c *gin.Context) string {
	slug, _ := c.Get(CtxTenantSlug)
	s, _ := slug.(string)
	return s
}
