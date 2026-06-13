package monitoring

import (
	"time"

	"github.com/google/uuid"
)

// ─────────────────────────────────────────────────────────────
// MonitoringProfile  — check frequency template
// ─────────────────────────────────────────────────────────────

type MonitoringProfile struct {
	ID                    uuid.UUID
	TenantID              uuid.UUID
	Name                  string
	Description           string
	CheckIntervalSeconds  int
	RetryIntervalSeconds  int
	MaxRetries            int
	NotificationPeriod    string
	IsDefault             bool
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// ─────────────────────────────────────────────────────────────
// SSHKey  — stored encrypted, used for SSH-based monitoring
// ─────────────────────────────────────────────────────────────

type SSHKey struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	Name           string
	PublicKey      string
	PrivateKeyEnc  string // AES-256-GCM encrypted at rest
	Fingerprint    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// ─────────────────────────────────────────────────────────────
// MonitoredHost  — aggregate root
// ─────────────────────────────────────────────────────────────

type MonitoredHost struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	AssetID        *uuid.UUID
	AgentID        *uuid.UUID
	ProfileID      *uuid.UUID
	Name           string
	DisplayName    string
	IPAddress      string
	Hostname       string
	MonitoringType MonitoringType

	// SNMP
	SNMPVersion   string
	SNMPCommunity string
	SNMPPort      int
	SNMPTimeout   int
	SNMPV3Config  *SNMPV3Config

	// SSH
	SSHPort     int
	SSHUsername string
	SSHKeyID    *uuid.UUID

	// WMI
	WMIUsername string
	WMIDomain   string

	// State
	Status           HostStatus
	MaintenanceStart *time.Time
	MaintenanceEnd   *time.Time
	LastCheckAt      *time.Time
	LastStatus       *HostCheckStatus
	UptimePercent    *float64

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type MonitoringType string

const (
	MonitoringTypeAgent MonitoringType = "agent"
	MonitoringTypeSNMP  MonitoringType = "snmp"
	MonitoringTypeWMI   MonitoringType = "wmi"
	MonitoringTypeICMP  MonitoringType = "icmp"
	MonitoringTypeSSH   MonitoringType = "ssh"
	MonitoringTypeAPI   MonitoringType = "api"
)

type HostStatus string

const (
	HostStatusActive      HostStatus = "active"
	HostStatusMaintenance HostStatus = "maintenance"
	HostStatusDisabled    HostStatus = "disabled"
	HostStatusPending     HostStatus = "pending"
)

type HostCheckStatus string

const (
	HostCheckStatusUp          HostCheckStatus = "up"
	HostCheckStatusDown        HostCheckStatus = "down"
	HostCheckStatusWarning     HostCheckStatus = "warning"
	HostCheckStatusUnknown     HostCheckStatus = "unknown"
	HostCheckStatusMaintenance HostCheckStatus = "maintenance"
)

type SNMPV3Config struct {
	Username  string `json:"username"`
	AuthProto string `json:"auth_proto"`  // MD5, SHA, SHA256, SHA512
	AuthKey   string `json:"auth_key"`    // encrypted
	PrivProto string `json:"priv_proto"`  // DES, AES, AES192, AES256
	PrivKey   string `json:"priv_key"`    // encrypted
	SecLevel  string `json:"sec_level"`   // noAuthNoPriv, authNoPriv, authPriv
}

func (h *MonitoredHost) IsInMaintenance() bool {
	now := time.Now()
	if h.MaintenanceStart == nil || h.MaintenanceEnd == nil {
		return h.Status == HostStatusMaintenance
	}
	return now.After(*h.MaintenanceStart) && now.Before(*h.MaintenanceEnd)
}

// ─────────────────────────────────────────────────────────────
// CheckDefinition  — what to measure on a host
// ─────────────────────────────────────────────────────────────

type CheckDefinition struct {
	ID                   uuid.UUID
	TenantID             uuid.UUID
	HostID               uuid.UUID
	Name                 string
	CheckType            CheckType
	Parameters           map[string]any
	Thresholds           CheckThresholds
	CheckIntervalSeconds *int // nil = inherit from profile
	IsEnabled            bool
	LastValue            *float64
	LastStatus           *string
	LastCheckedAt        *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type CheckType string

const (
	CheckTypeCPU       CheckType = "cpu"
	CheckTypeRAM       CheckType = "ram"
	CheckTypeDisk      CheckType = "disk"
	CheckTypeNetwork   CheckType = "network"
	CheckTypeService   CheckType = "service"
	CheckTypeProcess   CheckType = "process"
	CheckTypeICMP      CheckType = "icmp"
	CheckTypeHTTP      CheckType = "http"
	CheckTypeHTTPS     CheckType = "https"
	CheckTypeTCP       CheckType = "tcp"
	CheckTypeUDP       CheckType = "udp"
	CheckTypeDNS       CheckType = "dns"
	CheckTypeSMTP      CheckType = "smtp"
	CheckTypeSNMPOID   CheckType = "snmp_oid"
	CheckTypeCustom    CheckType = "custom"
)

type CheckThresholds struct {
	Warning  *float64 `json:"warning,omitempty"`
	Critical *float64 `json:"critical,omitempty"`
}

// ─────────────────────────────────────────────────────────────
// MetricPoint  — single time-series data point
// ─────────────────────────────────────────────────────────────

type MetricPoint struct {
	Time       time.Time
	TenantID   uuid.UUID
	HostID     uuid.UUID
	CheckID    uuid.UUID
	MetricName string
	Value      float64
	Labels     map[string]string
}

// ─────────────────────────────────────────────────────────────
// AvailabilityRecord  — state transition log for hosts
// ─────────────────────────────────────────────────────────────

type AvailabilityRecord struct {
	ID              uuid.UUID
	TenantID        uuid.UUID
	HostID          uuid.UUID
	Status          HostCheckStatus
	StartedAt       time.Time
	EndedAt         *time.Time
	DurationSeconds *int
	Reason          string
}

// ─────────────────────────────────────────────────────────────
// ServiceCheck  — application-level probes
// ─────────────────────────────────────────────────────────────

type ServiceCheck struct {
	ID                 uuid.UUID
	TenantID           uuid.UUID
	HostID             uuid.UUID
	ServiceName        string
	ServiceDisplayName string
	CheckType          string
	TargetHost         string
	TargetPort         *int
	ExpectedResponse   string
	ExpectedHTTPCode   *int
	TimeoutSeconds     int
	FollowRedirects    bool
	IsEnabled          bool
	LastStatus         *string
	LastCheckedAt      *time.Time
	LastResponseTimeMs *int
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
