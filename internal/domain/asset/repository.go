package asset

import (
	"context"

	"github.com/google/uuid"
	"github.com/infracore/infracore/internal/domain/shared"
)

// AssetFilter provides fine-grained filtering for asset queries.
type AssetFilter struct {
	SiteID       *uuid.UUID
	DepartmentID *uuid.UUID
	RackID       *uuid.UUID
	Category     *AssetCategory
	Status       *AssetStatus
	Search       string // full-text: name, serial, asset_tag
	Tags         []string
	HasWarranty  *bool
	WarrantyDaysRemaining *int // expiring within N days
}

type AssetRepository interface {
	Create(ctx context.Context, a *Asset) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*Asset, error)
	GetByAssetTag(ctx context.Context, tenantID uuid.UUID, tag string) (*Asset, error)
	GetBySerial(ctx context.Context, tenantID uuid.UUID, serial string) (*Asset, error)
	Update(ctx context.Context, a *Asset) error
	SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error
	List(ctx context.Context, tenantID uuid.UUID, filter AssetFilter, page shared.Page, sort shared.Sort) (shared.PageResult[*Asset], error)
	CountByCategory(ctx context.Context, tenantID uuid.UUID) (map[AssetCategory]int64, error)
	CountByStatus(ctx context.Context, tenantID uuid.UUID) (map[AssetStatus]int64, error)
	GetExpiringWarranties(ctx context.Context, tenantID uuid.UUID, withinDays int) ([]*Asset, error)
}

type ManufacturerRepository interface {
	Create(ctx context.Context, m *Manufacturer) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*Manufacturer, error)
	Update(ctx context.Context, m *Manufacturer) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
	List(ctx context.Context, tenantID uuid.UUID, page shared.Page) (shared.PageResult[*Manufacturer], error)
}

type DeviceTypeRepository interface {
	Create(ctx context.Context, dt *DeviceType) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*DeviceType, error)
	Update(ctx context.Context, dt *DeviceType) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
	List(ctx context.Context, tenantID uuid.UUID, category *AssetCategory, page shared.Page) (shared.PageResult[*DeviceType], error)
}

type NetworkInterfaceRepository interface {
	Create(ctx context.Context, ni *NetworkInterface) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*NetworkInterface, error)
	Update(ctx context.Context, ni *NetworkInterface) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
	ListByAsset(ctx context.Context, tenantID, assetID uuid.UUID) ([]*NetworkInterface, error)
	GetByMAC(ctx context.Context, tenantID uuid.UUID, mac string) (*NetworkInterface, error)
}

type VirtualMachineRepository interface {
	Create(ctx context.Context, vm *VirtualMachine) error
	GetByAssetID(ctx context.Context, tenantID, assetID uuid.UUID) (*VirtualMachine, error)
	GetByExternalID(ctx context.Context, tenantID uuid.UUID, hypervisor HypervisorType, externalID string) (*VirtualMachine, error)
	Update(ctx context.Context, vm *VirtualMachine) error
	ListByHost(ctx context.Context, tenantID, hostAssetID uuid.UUID) ([]*VirtualMachine, error)
}
