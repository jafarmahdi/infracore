package agent

import (
	"time"

	"github.com/google/uuid"
)

// ─────────────────────────────────────────────────────────────
// AgentGroup
// ─────────────────────────────────────────────────────────────

type AgentGroup struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	Name        string
	Description string
	IsDefault   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ─────────────────────────────────────────────────────────────
// Agent  — aggregate root
// ─────────────────────────────────────────────────────────────

type Agent struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	SiteID         *uuid.UUID
	GroupID        *uuid.UUID
	Name           string
	Description    string
	Hostname       string
	IPAddress      string
	OSType         AgentOSType
	OSVersion      string
	Arch           AgentArch
	Version        string
	APIKeyHash     string // SHA-256 of the actual API key
	Status         AgentStatus
	LastHeartbeatAt *time.Time
	LastSeenIP     string
	Capabilities   []AgentCapability
	Configuration  map[string]any
	RegisteredAt   time.Time
	UpdatedAt      time.Time
}

type AgentOSType string

const (
	AgentOSLinux   AgentOSType = "linux"
	AgentOSWindows AgentOSType = "windows"
	AgentOSDarwin  AgentOSType = "darwin"
)

type AgentArch string

const (
	AgentArchAMD64 AgentArch = "amd64"
	AgentArchARM64 AgentArch = "arm64"
	AgentArchARM   AgentArch = "arm"
	AgentArch386   AgentArch = "386"
)

type AgentStatus string

const (
	AgentStatusOnline      AgentStatus = "online"
	AgentStatusOffline     AgentStatus = "offline"
	AgentStatusError       AgentStatus = "error"
	AgentStatusUpdating    AgentStatus = "updating"
	AgentStatusMaintenance AgentStatus = "maintenance"
)

type AgentCapability string

const (
	CapabilityAgent       AgentCapability = "agent"     // native agent monitoring
	CapabilitySNMP        AgentCapability = "snmp"      // SNMP polling
	CapabilityWMI         AgentCapability = "wmi"       // WMI queries (Windows only)
	CapabilitySSH         AgentCapability = "ssh"       // SSH-based checks
	CapabilityICMP        AgentCapability = "icmp"      // Ping checks
	CapabilityDiscovery   AgentCapability = "discovery" // Network discovery
	CapabilityVMware      AgentCapability = "vmware"    // VMware API
	CapabilityHyperV      AgentCapability = "hyperv"    // Hyper-V API
)

func (a *Agent) IsOnline() bool {
	return a.Status == AgentStatusOnline
}

func (a *Agent) HeartbeatAge() *time.Duration {
	if a.LastHeartbeatAt == nil {
		return nil
	}
	age := time.Since(*a.LastHeartbeatAt)
	return &age
}

func (a *Agent) HasCapability(c AgentCapability) bool {
	for _, cap := range a.Capabilities {
		if cap == c {
			return true
		}
	}
	return false
}

// ─────────────────────────────────────────────────────────────
// AgentHealthMetric  — real-time agent performance data
// ─────────────────────────────────────────────────────────────

type AgentHealthMetric struct {
	Time             time.Time
	AgentID          uuid.UUID
	TenantID         uuid.UUID
	CPUPercent       float64
	RAMPercent       float64
	Goroutines       int
	ChecksPerMinute  int
	ErrorsPerMinute  int
	QueueDepth       int
	LatencyMs        int
}

// ─────────────────────────────────────────────────────────────
// AgentTask  — command queued for agent execution
// ─────────────────────────────────────────────────────────────

type AgentTask struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	AgentID     uuid.UUID
	TaskType    AgentTaskType
	Payload     map[string]any
	Status      AgentTaskStatus
	Result      map[string]any
	ErrorMsg    string
	CreatedAt   time.Time
	DeliveredAt *time.Time
	CompletedAt *time.Time
	ExpiresAt   time.Time
}

type AgentTaskType string

const (
	TaskTypeUpdateAgent    AgentTaskType = "update_agent"
	TaskTypeRunDiscovery   AgentTaskType = "run_discovery"
	TaskTypeReloadConfig   AgentTaskType = "reload_config"
	TaskTypeRunCheck       AgentTaskType = "run_check"
	TaskTypeCollectLogs    AgentTaskType = "collect_logs"
)

type AgentTaskStatus string

const (
	AgentTaskPending   AgentTaskStatus = "pending"
	AgentTaskDelivered AgentTaskStatus = "delivered"
	AgentTaskRunning   AgentTaskStatus = "running"
	AgentTaskCompleted AgentTaskStatus = "completed"
	AgentTaskFailed    AgentTaskStatus = "failed"
	AgentTaskExpired   AgentTaskStatus = "expired"
)

func (t *AgentTask) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// ─────────────────────────────────────────────────────────────
// AgentUpdateChannel  — version distribution channels
// ─────────────────────────────────────────────────────────────

type AgentUpdateChannel struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	ChannelName    string // stable, beta, edge
	OSType         string
	Arch           string
	Version        string
	DownloadURL    string
	ChecksumSHA256 string
	ReleaseNotes   string
	IsActive       bool
	CreatedAt      time.Time
}
