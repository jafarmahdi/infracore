package asset

import (
	"time"

	"github.com/google/uuid"
)

// ─────────────────────────────────────────────────────────────
// Manufacturer  — entity
// ─────────────────────────────────────────────────────────────

type Manufacturer struct {
	ID           uuid.UUID
	TenantID     uuid.UUID
	Name         string
	Description  string
	WebsiteURL   string
	SupportURL   string
	SupportPhone string
	SupportEmail string
	IsGlobal     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ─────────────────────────────────────────────────────────────
// DeviceType  — hardware model template
// ─────────────────────────────────────────────────────────────

type DeviceType struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	ManufacturerID uuid.UUID
	Name           string
	Model          string
	PartNumber     string
	Category       AssetCategory
	FormFactor     FormFactor
	RackUnits      *int
	WeightKg       *float64
	FrontImageURL  string
	RearImageURL   string
	Specs          map[string]any
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// ─────────────────────────────────────────────────────────────
// Asset  — aggregate root (unified model for all physical assets)
// ─────────────────────────────────────────────────────────────

type Asset struct {
	ID               uuid.UUID
	TenantID         uuid.UUID
	SiteID           *uuid.UUID
	DepartmentID     *uuid.UUID
	DeviceTypeID     *uuid.UUID
	RackID           *uuid.UUID
	RackPositionStart *int
	RackPositionEnd   *int
	RackFace         *string // "front" | "rear"

	AssetTag     string
	SerialNumber string
	Name         string
	Category     AssetCategory
	Status       AssetStatus
	Role         string // e.g. "web-server", "core-switch"

	// Denormalized for fast list rendering (avoids joins)
	ManufacturerName string
	ModelName        string

	// Connectivity
	PrimaryIP  string // CIDR notation
	OOBIP      string // Out-of-band management IP
	MACAddress string

	// Operating System
	OSType    string
	OSVersion string

	// Hardware specs
	CPUCount   *int
	CPUModel   string
	RAMGB      *int
	StorageGB  *int

	// Lifecycle
	PurchaseDate     *time.Time
	PurchaseCost     *float64
	PurchaseCurrency string
	PurchaseOrder    string
	WarrantyExpiry   *time.Time
	EOLDate          *time.Time

	// Ownership
	AssignedToUserID *uuid.UUID
	ManagedByUserID  *uuid.UUID

	Notes        string
	CustomFields map[string]any
	Tags         []string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (a *Asset) IsRackMounted() bool {
	return a.RackID != nil && a.RackPositionStart != nil
}

func (a *Asset) IsDecommissioned() bool {
	return a.Status == AssetStatusDecommissioned
}

func (a *Asset) WarrantyExpired() bool {
	if a.WarrantyExpiry == nil {
		return false
	}
	return time.Now().After(*a.WarrantyExpiry)
}

func (a *Asset) DaysUntilWarrantyExpiry() *int {
	if a.WarrantyExpiry == nil {
		return nil
	}
	days := int(time.Until(*a.WarrantyExpiry).Hours() / 24)
	return &days
}

// ─────────────────────────────────────────────────────────────
// VirtualMachine  — extends Asset with hypervisor metadata
// ─────────────────────────────────────────────────────────────

type VirtualMachine struct {
	ID             uuid.UUID
	AssetID        uuid.UUID
	TenantID       uuid.UUID
	HostAssetID    *uuid.UUID
	HypervisorType HypervisorType
	VMIDExternal   string
	VCPUCount      *int
	VRAMGB         *int
	StorageGB      *int
	ClusterName    string
	DatastoreName  string
	TemplateName   string
	PowerState     VMPowerState
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type HypervisorType string

const (
	HypervisorVMware  HypervisorType = "vmware"
	HypervisorHyperV  HypervisorType = "hyperv"
	HypervisorProxmox HypervisorType = "proxmox"
	HypervisorKVM     HypervisorType = "kvm"
	HypervisorXen     HypervisorType = "xen"
	HypervisorOther   HypervisorType = "other"
)

type VMPowerState string

const (
	VMPowerRunning   VMPowerState = "running"
	VMPowerStopped   VMPowerState = "stopped"
	VMPowerSuspended VMPowerState = "suspended"
	VMPowerPaused    VMPowerState = "paused"
	VMPowerUnknown   VMPowerState = "unknown"
)

// ─────────────────────────────────────────────────────────────
// NetworkInterface  — child of Asset
// ─────────────────────────────────────────────────────────────

type NetworkInterface struct {
	ID              uuid.UUID
	TenantID        uuid.UUID
	AssetID         uuid.UUID
	Name            string
	Description     string
	InterfaceType   string
	SpeedMbps       *int
	Duplex          string
	MACAddress      string
	MTU             int
	IsManagement    bool
	IsEnabled       bool
	VLANMode        string
	UntaggedVLANID  *uuid.UUID
	TaggedVLANIDs   []uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ─────────────────────────────────────────────────────────────
// AssetRelationship  — parent/child, cluster, HA pair
// ─────────────────────────────────────────────────────────────

type AssetRelationship struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	ParentAssetID  uuid.UUID
	ChildAssetID   uuid.UUID
	Relationship   RelationshipType
	CreatedAt      time.Time
}

type RelationshipType string

const (
	RelationshipParentChild RelationshipType = "parent_child"
	RelationshipCluster     RelationshipType = "cluster"
	RelationshipStack       RelationshipType = "stack"
	RelationshipHAPair      RelationshipType = "ha_pair"
)

// ─────────────────────────────────────────────────────────────
// Value types
// ─────────────────────────────────────────────────────────────

type AssetCategory string

const (
	AssetCategoryServer      AssetCategory = "server"
	AssetCategorySwitch      AssetCategory = "switch"
	AssetCategoryRouter      AssetCategory = "router"
	AssetCategoryFirewall    AssetCategory = "firewall"
	AssetCategoryUPS         AssetCategory = "ups"
	AssetCategoryPrinter     AssetCategory = "printer"
	AssetCategoryAccessPoint AssetCategory = "access_point"
	AssetCategoryStorage     AssetCategory = "storage"
	AssetCategoryVM          AssetCategory = "vm"
	AssetCategoryOther       AssetCategory = "other"
)

type AssetStatus string

const (
	AssetStatusActive         AssetStatus = "active"
	AssetStatusSpare          AssetStatus = "spare"
	AssetStatusOffline        AssetStatus = "offline"
	AssetStatusFailed         AssetStatus = "failed"
	AssetStatusDecommissioned AssetStatus = "decommissioned"
	AssetStatusPlanned        AssetStatus = "planned"
	AssetStatusInRepair       AssetStatus = "in_repair"
	AssetStatusInTransit      AssetStatus = "in_transit"
)

type FormFactor string

const (
	FormFactor1U        FormFactor = "1u"
	FormFactor2U        FormFactor = "2u"
	FormFactor4U        FormFactor = "4u"
	FormFactorBlade     FormFactor = "blade"
	FormFactorTower     FormFactor = "tower"
	FormFactorDesktop   FormFactor = "desktop"
	FormFactorModular   FormFactor = "modular"
	FormFactorWallMount FormFactor = "wall-mount"
	FormFactorOther     FormFactor = "other"
)
