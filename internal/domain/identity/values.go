package identity

// Address is a value object representing a postal address.
type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"` // ISO 3166-1 alpha-2
	Postal  string `json:"postal"`
}

// Coordinates is a value object for GPS location.
type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Contact holds basic contact information.
type Contact struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// TenantPlan defines the billing tier for a tenant.
type TenantPlan string

const (
	TenantPlanFree       TenantPlan = "free"
	TenantPlanPro        TenantPlan = "pro"
	TenantPlanEnterprise TenantPlan = "enterprise"
)

// TenantSettings holds tenant-level configuration.
type TenantSettings struct {
	AllowSSOLogin   bool           `json:"allow_sso_login"`
	DefaultTimeZone string         `json:"default_timezone"`
	PasswordPolicy  PasswordPolicy `json:"password_policy"`
	SessionTTLHours int            `json:"session_ttl_hours"`
	MFARequired     bool           `json:"mfa_required"`
}

// PasswordPolicy defines password strength requirements.
type PasswordPolicy struct {
	MinLength        int  `json:"min_length"`
	RequireUppercase bool `json:"require_uppercase"`
	RequireNumbers   bool `json:"require_numbers"`
	RequireSymbols   bool `json:"require_symbols"`
	ExpiryDays       int  `json:"expiry_days"`
	PreventReuse     int  `json:"prevent_reuse"` // disallow last N passwords
}

// UserPreferences stores per-user UI/UX settings.
type UserPreferences struct {
	Theme        string `json:"theme"`          // light, dark, system
	Language     string `json:"language"`       // ISO 639-1
	TimeZone     string `json:"timezone"`
	DateFormat   string `json:"date_format"`    // YYYY-MM-DD, DD/MM/YYYY, etc.
	TimeFormat   string `json:"time_format"`    // 24h, 12h
	ItemsPerPage int    `json:"items_per_page"`
	NotifyEmail  bool   `json:"notify_email"`
	NotifySMS    bool   `json:"notify_sms"`
}

// RoleSlug is a well-known built-in role identifier.
type RoleSlug string

const (
	RoleSlugSuperAdmin   RoleSlug = "super-admin"
	RoleSlugAdmin        RoleSlug = "admin"
	RoleSlugReadOnly     RoleSlug = "read-only"
	RoleSlugNetworkAdmin RoleSlug = "network-admin"
	RoleSlugAssetManager RoleSlug = "asset-manager"
	RoleSlugMonitor      RoleSlug = "monitor"
)
