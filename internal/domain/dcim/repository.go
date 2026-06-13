package dcim

import (
	"context"

	"github.com/google/uuid"
	"github.com/infracore/infracore/internal/domain/shared"
)

type DataCenterFilter struct {
	SiteID *uuid.UUID
	Status *DataCenterStatus
	Search string
}

type DataCenterRepository interface {
	Create(ctx context.Context, dc *DataCenter) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*DataCenter, error)
	Update(ctx context.Context, dc *DataCenter) error
	SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error
	List(ctx context.Context, tenantID uuid.UUID, filter DataCenterFilter, page shared.Page) (shared.PageResult[*DataCenter], error)
}

type RoomRepository interface {
	Create(ctx context.Context, r *Room) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*Room, error)
	Update(ctx context.Context, r *Room) error
	SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error
	ListByDataCenter(ctx context.Context, tenantID, dataCenterID uuid.UUID, page shared.Page) (shared.PageResult[*Room], error)
}

type RackFilter struct {
	SiteID *uuid.UUID
	RoomID *uuid.UUID
	Status *RackStatus
	Search string
}

type RackRepository interface {
	Create(ctx context.Context, r *Rack) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*Rack, error)
	Update(ctx context.Context, r *Rack) error
	SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error
	List(ctx context.Context, tenantID uuid.UUID, filter RackFilter, page shared.Page) (shared.PageResult[*Rack], error)
	GetOccupiedUnits(ctx context.Context, tenantID, rackID uuid.UUID) ([]RackUnitOccupancy, error)
}

// RackUnitOccupancy describes which asset occupies a range of rack units.
type RackUnitOccupancy struct {
	AssetID   uuid.UUID
	AssetName string
	StartUnit int
	EndUnit   int
	Face      string // front, rear
}

type PowerFeedRepository interface {
	Create(ctx context.Context, pf *PowerFeed) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*PowerFeed, error)
	Update(ctx context.Context, pf *PowerFeed) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
	ListByRack(ctx context.Context, tenantID, rackID uuid.UUID) ([]*PowerFeed, error)
}

type PDURepository interface {
	Create(ctx context.Context, p *PDU) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*PDU, error)
	Update(ctx context.Context, p *PDU) error
	SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error
	ListByRack(ctx context.Context, tenantID, rackID uuid.UUID) ([]*PDU, error)
}

type PatchPanelRepository interface {
	Create(ctx context.Context, p *PatchPanel) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*PatchPanel, error)
	Update(ctx context.Context, p *PatchPanel) error
	SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error
	ListByRack(ctx context.Context, tenantID, rackID uuid.UUID) ([]*PatchPanel, error)
}

type CableRepository interface {
	Create(ctx context.Context, c *Cable) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*Cable, error)
	Update(ctx context.Context, c *Cable) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
	List(ctx context.Context, tenantID uuid.UUID, page shared.Page) (shared.PageResult[*Cable], error)
	// GetByTermination finds cables connected to a specific interface/port.
	GetByTermination(ctx context.Context, terminationType string, terminationID uuid.UUID) ([]*Cable, error)
}
