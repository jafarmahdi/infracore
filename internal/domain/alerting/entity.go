package alerting

import (
	"time"

	"github.com/google/uuid"
)

// ─────────────────────────────────────────────────────────────
// NotificationChannel
// ─────────────────────────────────────────────────────────────

type NotificationChannel struct {
	ID            uuid.UUID
	TenantID      uuid.UUID
	Name          string
	ChannelType   ChannelType
	Configuration map[string]any // encrypted at rest
	IsEnabled     bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ChannelType string

const (
	ChannelTypeEmail      ChannelType = "email"
	ChannelTypeSMS        ChannelType = "sms"
	ChannelTypeWhatsApp   ChannelType = "whatsapp"
	ChannelTypeWebhook    ChannelType = "webhook"
	ChannelTypeSlack      ChannelType = "slack"
	ChannelTypeTeams      ChannelType = "teams"
	ChannelTypePagerDuty  ChannelType = "pagerduty"
	ChannelTypeTelegram   ChannelType = "telegram"
)

// ─────────────────────────────────────────────────────────────
// EscalationPolicy
// ─────────────────────────────────────────────────────────────

type EscalationPolicy struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	Name        string
	Description string
	Steps       []EscalationStep
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type EscalationStep struct {
	ID                    uuid.UUID
	PolicyID              uuid.UUID
	StepOrder             int
	EscalateAfterMinutes  int
	ChannelIDs            []uuid.UUID
	UserIDs               []uuid.UUID
	CreatedAt             time.Time
}

// ─────────────────────────────────────────────────────────────
// AlertRule  — aggregate root
// ─────────────────────────────────────────────────────────────

type AlertRule struct {
	ID                       uuid.UUID
	TenantID                 uuid.UUID
	EscalationPolicyID       *uuid.UUID
	Name                     string
	Description              string
	RuleType                 AlertRuleType
	ScopeType                AlertScopeType
	ScopeIDs                 []uuid.UUID
	MetricName               string
	CheckType                string
	LabelFilters             map[string]string
	Condition                AlertCondition
	ThresholdWarning         *float64
	ThresholdCritical        *float64
	ThresholdEmergency       *float64
	EvaluationIntervalSecs   int
	EvaluationPeriods        int // consecutive breaching periods to fire
	AbsenceWindowSeconds     *int
	Severity                 AlertSeverity
	IsEnabled                bool
	ChannelIDs               []uuid.UUID
	CreatedAt                time.Time
	UpdatedAt                time.Time
}

type AlertRuleType string

const (
	AlertRuleThreshold AlertRuleType = "threshold"
	AlertRuleAnomaly   AlertRuleType = "anomaly"
	AlertRuleAbsence   AlertRuleType = "absence"
	AlertRuleChange    AlertRuleType = "change"
	AlertRuleComposite AlertRuleType = "composite"
)

type AlertScopeType string

const (
	AlertScopeHost       AlertScopeType = "host"
	AlertScopeAgentGroup AlertScopeType = "agent_group"
	AlertScopeSite       AlertScopeType = "site"
	AlertScopeTag        AlertScopeType = "tag"
	AlertScopeAll        AlertScopeType = "all"
)

type AlertCondition string

const (
	AlertConditionGT  AlertCondition = "gt"
	AlertConditionLT  AlertCondition = "lt"
	AlertConditionGTE AlertCondition = "gte"
	AlertConditionLTE AlertCondition = "lte"
	AlertConditionEQ  AlertCondition = "eq"
	AlertConditionNEQ AlertCondition = "neq"
)

type AlertSeverity string

const (
	AlertSeverityInfo      AlertSeverity = "info"
	AlertSeverityWarning   AlertSeverity = "warning"
	AlertSeverityCritical  AlertSeverity = "critical"
	AlertSeverityEmergency AlertSeverity = "emergency"
)

// ─────────────────────────────────────────────────────────────
// AlertEvent  — a fired alert instance
// ─────────────────────────────────────────────────────────────

type AlertEvent struct {
	ID                    uuid.UUID
	TenantID              uuid.UUID
	RuleID                uuid.UUID
	HostID                *uuid.UUID
	CheckID               *uuid.UUID
	Severity              AlertSeverity
	Status                AlertEventStatus
	Title                 string
	Message               string
	MetricName            string
	MetricValue           *float64
	ThresholdValue        *float64
	Labels                map[string]string
	FiredAt               time.Time
	ResolvedAt            *time.Time
	AcknowledgedAt        *time.Time
	AcknowledgedBy        *uuid.UUID
	AcknowledgementNote   string
	SilencedUntil         *time.Time
	SilencedBy            *uuid.UUID
	EscalationPolicyID    *uuid.UUID
	CurrentEscalationStep int
	LastEscalatedAt       *time.Time
}

type AlertEventStatus string

const (
	AlertEventFiring       AlertEventStatus = "firing"
	AlertEventResolved     AlertEventStatus = "resolved"
	AlertEventAcknowledged AlertEventStatus = "acknowledged"
	AlertEventSilenced     AlertEventStatus = "silenced"
)

func (e *AlertEvent) IsOpen() bool {
	return e.Status == AlertEventFiring
}

func (e *AlertEvent) Duration() time.Duration {
	end := time.Now()
	if e.ResolvedAt != nil {
		end = *e.ResolvedAt
	}
	return end.Sub(e.FiredAt)
}

func (e *AlertEvent) Acknowledge(byUserID uuid.UUID, note string) {
	now := time.Now()
	e.AcknowledgedAt = &now
	e.AcknowledgedBy = &byUserID
	e.AcknowledgementNote = note
	e.Status = AlertEventAcknowledged
}

func (e *AlertEvent) Resolve() {
	now := time.Now()
	e.ResolvedAt = &now
	e.Status = AlertEventResolved
}

// ─────────────────────────────────────────────────────────────
// AlertNotification  — delivery attempt record
// ─────────────────────────────────────────────────────────────

type AlertNotification struct {
	ID           uuid.UUID
	TenantID     uuid.UUID
	EventID      uuid.UUID
	ChannelID    uuid.UUID
	Status       NotificationStatus
	SentAt       *time.Time
	ErrorMessage string
	RetryCount   int
	Payload      map[string]any
	CreatedAt    time.Time
}

type NotificationStatus string

const (
	NotificationPending NotificationStatus = "pending"
	NotificationSent    NotificationStatus = "sent"
	NotificationFailed  NotificationStatus = "failed"
	NotificationSkipped NotificationStatus = "skipped"
)

// ─────────────────────────────────────────────────────────────
// Silence  — suppresses alert firing during a window
// ─────────────────────────────────────────────────────────────

type Silence struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	CreatedBy uuid.UUID
	Matchers  SilenceMatchers
	StartsAt  time.Time
	EndsAt    time.Time
	Comment   string
	CreatedAt time.Time
}

type SilenceMatchers struct {
	HostID  *uuid.UUID        `json:"host_id,omitempty"`
	RuleID  *uuid.UUID        `json:"rule_id,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
}

func (s *Silence) IsActive() bool {
	now := time.Now()
	return now.After(s.StartsAt) && now.Before(s.EndsAt)
}
