package contract

import (
	"time"

	"github.com/google/uuid"
)

// Contract  — aggregate root
type Contract struct {
	ID                    uuid.UUID
	TenantID              uuid.UUID
	SiteID                *uuid.UUID
	DepartmentID          *uuid.UUID
	VendorID              *uuid.UUID
	ContractNumber        string
	Title                 string
	ContractType          ContractType
	Description           string
	StartDate             time.Time
	EndDate               *time.Time
	AutoRenew             bool
	RenewalNoticeDays     int
	Value                 *float64
	Currency              string
	BillingCycle          BillingCycle
	Status                ContractStatus
	SLAResponseHours      *int
	SLAResolutionHours    *int
	ContactName           string
	ContactEmail          string
	ContactPhone          string
	Notes                 string
	DocumentURLs          []string
	NotifyDaysBefore      []int
	Tags                  []string
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             *time.Time
}

type ContractType string

const (
	ContractTypeMaintenance ContractType = "maintenance"
	ContractTypeSupport     ContractType = "support"
	ContractTypeWarranty    ContractType = "warranty"
	ContractTypeService     ContractType = "service"
	ContractTypeLease       ContractType = "lease"
	ContractTypeNDA         ContractType = "nda"
	ContractTypeSLA         ContractType = "sla"
	ContractTypeOther       ContractType = "other"
)

type ContractStatus string

const (
	ContractStatusDraft      ContractStatus = "draft"
	ContractStatusActive     ContractStatus = "active"
	ContractStatusExpired    ContractStatus = "expired"
	ContractStatusCancelled  ContractStatus = "cancelled"
	ContractStatusRenewed    ContractStatus = "renewed"
	ContractStatusUnderReview ContractStatus = "under_review"
)

type BillingCycle string

const (
	BillingCycleMonthly   BillingCycle = "monthly"
	BillingCycleQuarterly BillingCycle = "quarterly"
	BillingCycleAnnually  BillingCycle = "annually"
	BillingCycleOneTime   BillingCycle = "one_time"
)

func (c *Contract) IsExpired() bool {
	if c.EndDate == nil {
		return false
	}
	return time.Now().After(*c.EndDate)
}

func (c *Contract) DaysUntilExpiry() *int {
	if c.EndDate == nil {
		return nil
	}
	days := int(time.Until(*c.EndDate).Hours() / 24)
	return &days
}

// Warranty  — dedicated warranty record linked to an asset
type Warranty struct {
	ID                uuid.UUID
	TenantID          uuid.UUID
	AssetID           uuid.UUID
	VendorID          *uuid.UUID
	ContractID        *uuid.UUID
	WarrantyType      WarrantyType
	StartDate         time.Time
	EndDate           time.Time
	Description       string
	SupportContact    string
	RMAProcess        string
	NotifyDaysBefore  []int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type WarrantyType string

const (
	WarrantyTypeManufacturer WarrantyType = "manufacturer"
	WarrantyTypeExtended     WarrantyType = "extended"
	WarrantyTypeAccidental   WarrantyType = "accidental"
	WarrantyTypeOnSite       WarrantyType = "on_site"
	WarrantyTypeNextDay      WarrantyType = "next_day"
)

func (w *Warranty) IsExpired() bool {
	return time.Now().After(w.EndDate)
}

func (w *Warranty) DaysRemaining() int {
	return int(time.Until(w.EndDate).Hours() / 24)
}
