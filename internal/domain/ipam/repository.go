package ipam

import (
	"context"
	"net"

	"github.com/google/uuid"
	"github.com/infracore/infracore/internal/domain/shared"
)

type VRFRepository interface {
	Create(ctx context.Context, v *VRF) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*VRF, error)
	GetByName(ctx context.Context, tenantID uuid.UUID, name string) (*VRF, error)
	Update(ctx context.Context, v *VRF) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
	List(ctx context.Context, tenantID uuid.UUID, page shared.Page) (shared.PageResult[*VRF], error)
}

type VLANFilter struct {
	SiteID      *uuid.UUID
	VLANGroupID *uuid.UUID
	VRFID       *uuid.UUID
	Status      *VLANStatus
	Role        string
	VIDMin      *int
	VIDMax      *int
}

type VLANRepository interface {
	Create(ctx context.Context, v *VLAN) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*VLAN, error)
	GetByVID(ctx context.Context, tenantID uuid.UUID, siteID *uuid.UUID, vid int) (*VLAN, error)
	Update(ctx context.Context, v *VLAN) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
	List(ctx context.Context, tenantID uuid.UUID, filter VLANFilter, page shared.Page) (shared.PageResult[*VLAN], error)
}

type PrefixFilter struct {
	SiteID  *uuid.UUID
	VRFID   *uuid.UUID
	VLANID  *uuid.UUID
	Status  *PrefixStatus
	Family  *int    // 4 or 6
	Within  *net.IPNet  // find prefixes contained within this supernet
	Search  string
}

type PrefixRepository interface {
	Create(ctx context.Context, p *Prefix) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*Prefix, error)
	GetByNetwork(ctx context.Context, tenantID uuid.UUID, vrfID *uuid.UUID, network *net.IPNet) (*Prefix, error)
	Update(ctx context.Context, p *Prefix) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
	List(ctx context.Context, tenantID uuid.UUID, filter PrefixFilter, page shared.Page) (shared.PageResult[*Prefix], error)
	GetContaining(ctx context.Context, tenantID uuid.UUID, ip net.IP) ([]*Prefix, error)
	GetAvailablePrefixes(ctx context.Context, tenantID, parentID uuid.UUID, prefixLength int) ([]*net.IPNet, error)
	GetUtilization(ctx context.Context, tenantID, prefixID uuid.UUID) (PrefixUtilization, error)
}

// PrefixUtilization summarizes IP usage within a prefix.
type PrefixUtilization struct {
	TotalIPs     int64
	AllocatedIPs int64
	AvailableIPs int64
	Percent      float64
}

type IPAddressFilter struct {
	VRFID              *uuid.UUID
	PrefixID           *uuid.UUID
	Status             *IPAddressStatus
	Family             *int
	AssignedObjectType string
	AssignedObjectID   *uuid.UUID
	DNSName            string
	Search             string
}

type IPAddressRepository interface {
	Create(ctx context.Context, ip *IPAddress) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*IPAddress, error)
	GetByAddress(ctx context.Context, tenantID uuid.UUID, vrfID *uuid.UUID, address *net.IPNet) (*IPAddress, error)
	Update(ctx context.Context, ip *IPAddress) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
	List(ctx context.Context, tenantID uuid.UUID, filter IPAddressFilter, page shared.Page) (shared.PageResult[*IPAddress], error)
	GetNextAvailable(ctx context.Context, tenantID, prefixID uuid.UUID) (*net.IPNet, error)
	Assign(ctx context.Context, ipID uuid.UUID, objectType string, objectID uuid.UUID) error
	Unassign(ctx context.Context, ipID uuid.UUID) error
}

type DHCPLeaseRepository interface {
	Create(ctx context.Context, l *DHCPLease) error
	GetByIP(ctx context.Context, tenantID uuid.UUID, ip net.IP) (*DHCPLease, error)
	GetByMAC(ctx context.Context, tenantID uuid.UUID, mac string) ([]*DHCPLease, error)
	Update(ctx context.Context, l *DHCPLease) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
	ListByPrefix(ctx context.Context, tenantID, prefixID uuid.UUID, page shared.Page) (shared.PageResult[*DHCPLease], error)
	DeleteExpired(ctx context.Context) (int64, error)
}

type DNSZoneRepository interface {
	Create(ctx context.Context, z *DNSZone) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*DNSZone, error)
	GetByName(ctx context.Context, tenantID uuid.UUID, name string) (*DNSZone, error)
	Update(ctx context.Context, z *DNSZone) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
	List(ctx context.Context, tenantID uuid.UUID, page shared.Page) (shared.PageResult[*DNSZone], error)
}

type DNSRecordRepository interface {
	Create(ctx context.Context, r *DNSRecord) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*DNSRecord, error)
	Update(ctx context.Context, r *DNSRecord) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
	ListByZone(ctx context.Context, tenantID, zoneID uuid.UUID, recordType *DNSRecordType, page shared.Page) (shared.PageResult[*DNSRecord], error)
	GetByName(ctx context.Context, tenantID, zoneID uuid.UUID, name string, rtype DNSRecordType) (*DNSRecord, error)
}
