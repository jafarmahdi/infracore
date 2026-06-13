package discovery

import (
	"net"
	"time"

	"github.com/google/uuid"
)

// DiscoveryJob  — aggregate root for a discovery scan
type DiscoveryJob struct {
	ID              uuid.UUID
	TenantID        uuid.UUID
	SiteID          *uuid.UUID
	AgentID         *uuid.UUID
	Name            string
	JobType         DiscoveryJobType
	Status          DiscoveryJobStatus
	ScheduleCron    string // empty = run once
	NextRunAt       *time.Time
	LastRunAt       *time.Time
	Configuration   map[string]any
	ProgressPercent int
	DiscoveredCount int
	ImportedCount   int
	ErrorMessage    string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type DiscoveryJobType string

const (
	DiscoveryTypeNetwork         DiscoveryJobType = "network"
	DiscoveryTypeSNMP            DiscoveryJobType = "snmp"
	DiscoveryTypeVMware          DiscoveryJobType = "vmware"
	DiscoveryTypeHyperV          DiscoveryJobType = "hyperv"
	DiscoveryTypeProxmox         DiscoveryJobType = "proxmox"
	DiscoveryTypeActiveDirectory DiscoveryJobType = "active_directory"
	DiscoveryTypeNmap            DiscoveryJobType = "nmap"
)

type DiscoveryJobStatus string

const (
	DiscoveryStatusPending   DiscoveryJobStatus = "pending"
	DiscoveryStatusQueued    DiscoveryJobStatus = "queued"
	DiscoveryStatusRunning   DiscoveryJobStatus = "running"
	DiscoveryStatusCompleted DiscoveryJobStatus = "completed"
	DiscoveryStatusFailed    DiscoveryJobStatus = "failed"
	DiscoveryStatusCancelled DiscoveryJobStatus = "cancelled"
)

func (j *DiscoveryJob) IsRunning() bool {
	return j.Status == DiscoveryStatusRunning
}

func (j *DiscoveryJob) IsScheduled() bool {
	return j.ScheduleCron != ""
}

// DiscoveryResult  — a single device found during discovery
type DiscoveryResult struct {
	ID             uuid.UUID
	JobID          uuid.UUID
	TenantID       uuid.UUID
	IPAddress      net.IP
	Hostname       string
	FQDN           string
	MACAddress     string
	DeviceType     string
	Vendor         string
	Model          string
	OSInfo         string
	OpenPorts      []DiscoveredPort
	SNMPData       SNMPDiscoveryData
	VMwareData     map[string]any
	ADData         map[string]any
	RawData        map[string]any
	Status         DiscoveryResultStatus
	MatchedAssetID *uuid.UUID
	CreatedAt      time.Time
}

type DiscoveredPort struct {
	Port    int    `json:"port"`
	Proto   string `json:"proto"`   // tcp, udp
	Service string `json:"service"` // ssh, http, snmp, etc.
}

type SNMPDiscoveryData struct {
	SysDescr    string            `json:"sys_descr"`
	SysName     string            `json:"sys_name"`
	SysOID      string            `json:"sys_oid"`
	SysLocation string            `json:"sys_location"`
	SysContact  string            `json:"sys_contact"`
	Interfaces  []SNMPInterface   `json:"interfaces"`
	CustomOIDs  map[string]string `json:"custom_oids"`
}

type SNMPInterface struct {
	Index       int    `json:"index"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        int    `json:"type"`
	Speed       int64  `json:"speed"`
	MACAddress  string `json:"mac_address"`
	AdminStatus string `json:"admin_status"`
	OperStatus  string `json:"oper_status"`
}

type DiscoveryResultStatus string

const (
	DiscoveryResultNew           DiscoveryResultStatus = "new"
	DiscoveryResultImported      DiscoveryResultStatus = "imported"
	DiscoveryResultIgnored       DiscoveryResultStatus = "ignored"
	DiscoveryResultDuplicate     DiscoveryResultStatus = "duplicate"
	DiscoveryResultPendingReview DiscoveryResultStatus = "pending_review"
)

// DiscoveryRule  — auto-import / auto-tag rule
type DiscoveryRule struct {
	ID              uuid.UUID
	TenantID        uuid.UUID
	Name            string
	Description     string
	IsEnabled       bool
	Priority        int
	MatchConditions DiscoveryMatchConditions
	Actions         DiscoveryActions
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type DiscoveryMatchConditions struct {
	IPRange         string            `json:"ip_range,omitempty"`
	Vendor          string            `json:"vendor,omitempty"`
	OpenPorts       []int             `json:"open_ports,omitempty"`
	HostnamePattern string            `json:"hostname_pattern,omitempty"`
	SNMPOIDMatches  map[string]string `json:"snmp_oid_matches,omitempty"`
}

type DiscoveryActions struct {
	AutoImport bool              `json:"auto_import"`
	Category   string            `json:"category,omitempty"`
	SiteID     *uuid.UUID        `json:"site_id,omitempty"`
	Role       string            `json:"role,omitempty"`
	Tags       []string          `json:"tags,omitempty"`
	CustomFields map[string]any  `json:"custom_fields,omitempty"`
}
