package monitoring

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	domain "github.com/infracore/infracore/internal/domain/monitoring"
)

type monitoredHostRow struct {
	ID             uuid.UUID      `db:"id"`
	TenantID       uuid.UUID      `db:"tenant_id"`
	AssetID        *uuid.UUID     `db:"asset_id"`
	AgentID        *uuid.UUID     `db:"agent_id"`
	ProfileID      *uuid.UUID     `db:"profile_id"`
	Name           string         `db:"name"`
	DisplayName    sql.NullString `db:"display_name"`
	IPAddress      sql.NullString `db:"ip_address"`
	Hostname       sql.NullString `db:"hostname"`
	MonitoringType string         `db:"monitoring_type"`
	SNMPVersion    sql.NullString `db:"snmp_version"`
	SNMPCommunity  sql.NullString `db:"snmp_community"`
	SNMPPort       sql.NullInt32  `db:"snmp_port"`
	SNMPTimeout    sql.NullInt32  `db:"snmp_timeout_secs"`
	SSHPort        sql.NullInt32  `db:"ssh_port"`
	SSHUsername    sql.NullString `db:"ssh_username"`
	SSHKeyID       *uuid.UUID     `db:"ssh_key_id"`
	WMIUsername    sql.NullString `db:"wmi_username"`
	WMIDomain      sql.NullString `db:"wmi_domain"`
	Status         string         `db:"status"`
	LastCheckAt    *time.Time     `db:"last_check_at"`
	LastStatus     sql.NullString `db:"last_status"`
	UptimePercent  *float64       `db:"uptime_percent"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at"`
}

func (r *monitoredHostRow) toDomain() *domain.MonitoredHost {
	h := &domain.MonitoredHost{
		ID:             r.ID,
		TenantID:       r.TenantID,
		AssetID:        r.AssetID,
		AgentID:        r.AgentID,
		ProfileID:      r.ProfileID,
		Name:           r.Name,
		MonitoringType: domain.MonitoringType(r.MonitoringType),
		Status:         domain.HostStatus(r.Status),
		LastCheckAt:    r.LastCheckAt,
		UptimePercent:  r.UptimePercent,
		SSHKeyID:       r.SSHKeyID,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
	if r.DisplayName.Valid {
		h.DisplayName = r.DisplayName.String
	}
	if r.IPAddress.Valid {
		h.IPAddress = r.IPAddress.String
	}
	if r.Hostname.Valid {
		h.Hostname = r.Hostname.String
	}
	if r.SNMPVersion.Valid {
		h.SNMPVersion = r.SNMPVersion.String
	}
	if r.SNMPCommunity.Valid {
		h.SNMPCommunity = r.SNMPCommunity.String
	}
	if r.SNMPPort.Valid {
		h.SNMPPort = int(r.SNMPPort.Int32)
	}
	if r.SNMPTimeout.Valid {
		h.SNMPTimeout = int(r.SNMPTimeout.Int32)
	}
	if r.SSHPort.Valid {
		h.SSHPort = int(r.SSHPort.Int32)
	}
	if r.SSHUsername.Valid {
		h.SSHUsername = r.SSHUsername.String
	}
	if r.WMIUsername.Valid {
		h.WMIUsername = r.WMIUsername.String
	}
	if r.WMIDomain.Valid {
		h.WMIDomain = r.WMIDomain.String
	}
	if r.LastStatus.Valid {
		ls := domain.HostCheckStatus(r.LastStatus.String)
		h.LastStatus = &ls
	}
	return h
}
