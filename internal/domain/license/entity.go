package license

import (
	"time"

	"github.com/google/uuid"
)

// SoftwareLicense  — aggregate root
type SoftwareLicense struct {
	ID                    uuid.UUID
	TenantID              uuid.UUID
	SiteID                *uuid.UUID
	DepartmentID          *uuid.UUID
	VendorID              *uuid.UUID
	Name                  string
	ProductName           string
	Version               string
	LicenseType           LicenseType
	LicenseKey            string // encrypted at rest
	LicenseCount          int
	SeatsUsed             int
	PurchaseDate          *time.Time
	PurchaseCost          *float64
	PurchaseCurrency      string
	PurchaseOrder         string
	InvoiceNumber         string
	ExpirationDate        *time.Time
	SupportExpirationDate *time.Time
	RenewalCost           *float64
	Notes                 string
	Status                LicenseStatus
	NotifyDaysBefore      []int  // [30, 14, 7]
	DocumentURLs          []string
	Tags                  []string
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             *time.Time
}

type LicenseType string

const (
	LicenseTypePerpetual    LicenseType = "perpetual"
	LicenseTypeSubscription LicenseType = "subscription"
	LicenseTypeConcurrent   LicenseType = "concurrent"
	LicenseTypePerSeat      LicenseType = "per_seat"
	LicenseTypePerDevice    LicenseType = "per_device"
	LicenseTypePerCore      LicenseType = "per_core"
	LicenseTypePerServer    LicenseType = "per_server"
	LicenseTypeOpenSource   LicenseType = "open_source"
	LicenseTypeFreeware     LicenseType = "freeware"
	LicenseTypeTrial        LicenseType = "trial"
)

type LicenseStatus string

const (
	LicenseStatusActive   LicenseStatus = "active"
	LicenseStatusExpired  LicenseStatus = "expired"
	LicenseStatusCancelled LicenseStatus = "cancelled"
	LicenseStatusPending  LicenseStatus = "pending"
)

func (l *SoftwareLicense) AvailableSeats() int {
	return l.LicenseCount - l.SeatsUsed
}

func (l *SoftwareLicense) IsExpired() bool {
	if l.ExpirationDate == nil {
		return false
	}
	return time.Now().After(*l.ExpirationDate)
}

func (l *SoftwareLicense) DaysUntilExpiry() *int {
	if l.ExpirationDate == nil {
		return nil
	}
	days := int(time.Until(*l.ExpirationDate).Hours() / 24)
	return &days
}

func (l *SoftwareLicense) IsExhausted() bool {
	return l.SeatsUsed >= l.LicenseCount
}

// LicenseAssignment  — assignment of a license seat to an asset or user
type LicenseAssignment struct {
	ID                uuid.UUID
	TenantID          uuid.UUID
	LicenseID         uuid.UUID
	AssignedToType    string // "asset" | "user"
	AssignedToID      uuid.UUID
	LicenseKeySegment string
	AssignedAt        time.Time
	AssignedBy        uuid.UUID
	Notes             string
}
