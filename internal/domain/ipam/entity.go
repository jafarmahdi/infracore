package ipam

import (
	"net"
	"time"

	"github.com/google/uuid"
)

// ─────────────────────────────────────────────────────────────
// VRF  — Virtual Routing and Forwarding instance
// ─────────────────────────────────────────────────────────────

type VRF struct {
	ID            uuid.UUID
	TenantID      uuid.UUID
	Name          string
	RD            string // Route Distinguisher e.g. "65000:1"
	Description   string
	EnforceUnique bool // prevent duplicate IPs within this VRF
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// ─────────────────────────────────────────────────────────────
// VLANGroup
// ─────────────────────────────────────────────────────────────

type VLANGroup struct {
	ID           uuid.UUID
	TenantID     uuid.UUID
	SiteID       *uuid.UUID
	Name         string
	Slug         string
	VLANIDRanges []VLANIDRange
	Description  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type VLANIDRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// ─────────────────────────────────────────────────────────────
// VLAN
// ─────────────────────────────────────────────────────────────

type VLAN struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	SiteID      *uuid.UUID
	VLANGroupID *uuid.UUID
	VRFID       *uuid.UUID
	VID         int // 1–4094
	Name        string
	Description string
	Status      VLANStatus
	Role        string // data, voice, management, storage, native, quarantine
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type VLANStatus string

const (
	VLANStatusActive     VLANStatus = "active"
	VLANStatusReserved   VLANStatus = "reserved"
	VLANStatusDeprecated VLANStatus = "deprecated"
)

// ─────────────────────────────────────────────────────────────
// Prefix  — IP network / subnet
// ─────────────────────────────────────────────────────────────

type Prefix struct {
	ID           uuid.UUID
	TenantID     uuid.UUID
	SiteID       *uuid.UUID
	VRFID        *uuid.UUID
	VLANID       *uuid.UUID
	Network      *net.IPNet // e.g. 192.168.1.0/24
	Status       PrefixStatus
	Role         string // loopback, link, management, container, pool
	Description  string
	IsPool       bool // addresses can be auto-allocated from this prefix
	MarkUtilized bool
	Tags         []string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type PrefixStatus string

const (
	PrefixStatusActive     PrefixStatus = "active"
	PrefixStatusReserved   PrefixStatus = "reserved"
	PrefixStatusDeprecated PrefixStatus = "deprecated"
	PrefixStatusContainer  PrefixStatus = "container" // supernet, not directly allocated
)

func (p *Prefix) Family() int {
	if p.Network == nil {
		return 0
	}
	if p.Network.IP.To4() != nil {
		return 4
	}
	return 6
}

func (p *Prefix) PrefixLength() int {
	if p.Network == nil {
		return 0
	}
	ones, _ := p.Network.Mask.Size()
	return ones
}

// ─────────────────────────────────────────────────────────────
// IPAddress
// ─────────────────────────────────────────────────────────────

type IPAddress struct {
	ID                 uuid.UUID
	TenantID           uuid.UUID
	VRFID              *uuid.UUID
	PrefixID           *uuid.UUID
	Address            *net.IPNet // includes prefix length: 192.168.1.10/24
	Status             IPAddressStatus
	Role               IPAddressRole
	AssignedObjectType string    // "network_interface" | "virtual_machine" | "asset"
	AssignedObjectID   *uuid.UUID
	DNSName            string
	Description        string
	NATInsideID        *uuid.UUID // points to another IPAddress
	Tags               []string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type IPAddressStatus string

const (
	IPStatusActive     IPAddressStatus = "active"
	IPStatusReserved   IPAddressStatus = "reserved"
	IPStatusDeprecated IPAddressStatus = "deprecated"
	IPStatusDHCP       IPAddressStatus = "dhcp"
	IPStatusSLAAC      IPAddressStatus = "slaac"
)

type IPAddressRole string

const (
	IPRoleLoopback   IPAddressRole = "loopback"
	IPRoleSecondary  IPAddressRole = "secondary"
	IPRoleAnycast    IPAddressRole = "anycast"
	IPRoleVirtual    IPAddressRole = "virtual"
	IPRoleVIP        IPAddressRole = "vip"
	IPRoleVRRP       IPAddressRole = "vrrp"
	IPRoleHSRP       IPAddressRole = "hsrp"
	IPRoleGLBP       IPAddressRole = "glbp"
)

func (ip *IPAddress) Family() int {
	if ip.Address == nil {
		return 0
	}
	if ip.Address.IP.To4() != nil {
		return 4
	}
	return 6
}

func (ip *IPAddress) IsAssigned() bool {
	return ip.AssignedObjectID != nil
}

// ─────────────────────────────────────────────────────────────
// DHCPLease
// ─────────────────────────────────────────────────────────────

type DHCPLease struct {
	ID         uuid.UUID
	TenantID   uuid.UUID
	PrefixID   *uuid.UUID
	IPAddress  net.IP
	MACAddress string
	Hostname   string
	ClientID   string
	AssetID    *uuid.UUID
	LeaseStart *time.Time
	LeaseEnd   *time.Time
	IsStatic   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (l *DHCPLease) IsExpired() bool {
	if l.LeaseEnd == nil {
		return false
	}
	return time.Now().After(*l.LeaseEnd)
}

// ─────────────────────────────────────────────────────────────
// DNSZone
// ─────────────────────────────────────────────────────────────

type DNSZone struct {
	ID           uuid.UUID
	TenantID     uuid.UUID
	Name         string
	Description  string
	ZoneType     DNSZoneType
	IsActive     bool
	SOAPrimaryNS string
	SOAEmail     string
	SOASerial    int
	SOARefresh   int
	SOARetry     int
	SOAExpire    int
	SOAMinimum   int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type DNSZoneType string

const (
	DNSZoneTypeForward DNSZoneType = "forward"
	DNSZoneTypeReverse DNSZoneType = "reverse"
)

// ─────────────────────────────────────────────────────────────
// DNSRecord
// ─────────────────────────────────────────────────────────────

type DNSRecord struct {
	ID            uuid.UUID
	TenantID      uuid.UUID
	ZoneID        uuid.UUID
	Name          string // relative label; "@" for zone apex
	RecordType    DNSRecordType
	Value         string
	TTL           int
	Priority      *int // MX, SRV
	IPAddressID   *uuid.UUID
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type DNSRecordType string

const (
	DNSRecordA     DNSRecordType = "A"
	DNSRecordAAAA  DNSRecordType = "AAAA"
	DNSRecordCNAME DNSRecordType = "CNAME"
	DNSRecordMX    DNSRecordType = "MX"
	DNSRecordTXT   DNSRecordType = "TXT"
	DNSRecordPTR   DNSRecordType = "PTR"
	DNSRecordSRV   DNSRecordType = "SRV"
	DNSRecordNS    DNSRecordType = "NS"
	DNSRecordCAA   DNSRecordType = "CAA"
	DNSRecordSOA   DNSRecordType = "SOA"
)
