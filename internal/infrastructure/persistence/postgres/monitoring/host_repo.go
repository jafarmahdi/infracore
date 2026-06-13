package monitoring

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	domain "github.com/infracore/infracore/internal/domain/monitoring"
	"github.com/infracore/infracore/internal/domain/shared"
)

type HostRepository struct {
	db *sqlx.DB
}

func NewHostRepository(db *sqlx.DB) *HostRepository {
	return &HostRepository{db: db}
}

const hostCols = `
	id, tenant_id, asset_id, agent_id, profile_id,
	name, display_name, ip_address::text AS ip_address, hostname, monitoring_type,
	snmp_version, snmp_community, snmp_port, snmp_timeout_secs,
	ssh_port, ssh_username, ssh_key_id, wmi_username, wmi_domain,
	status, last_check_at, last_status, uptime_percent,
	created_at, updated_at`

func (r *HostRepository) Create(ctx context.Context, h *domain.MonitoredHost) error {
	h.ID = uuid.New()

	type params struct {
		ID             uuid.UUID `db:"id"`
		TenantID       uuid.UUID `db:"tenant_id"`
		Name           string    `db:"name"`
		DisplayName    *string   `db:"display_name"`
		IPAddress      *string   `db:"ip_address"`
		Hostname       *string   `db:"hostname"`
		MonitoringType string    `db:"monitoring_type"`
		SNMPVersion    *string   `db:"snmp_version"`
		SNMPCommunity  *string   `db:"snmp_community"`
		SNMPPort       *int      `db:"snmp_port"`
		SNMPTimeout    *int      `db:"snmp_timeout_secs"`
		SSHPort        *int      `db:"ssh_port"`
		SSHUsername    *string   `db:"ssh_username"`
		Status         string    `db:"status"`
	}

	p := params{
		ID:             h.ID,
		TenantID:       h.TenantID,
		Name:           h.Name,
		MonitoringType: string(h.MonitoringType),
		Status:         string(h.Status),
	}
	if h.DisplayName != "" {
		p.DisplayName = &h.DisplayName
	}
	if h.IPAddress != "" {
		p.IPAddress = &h.IPAddress
	}
	if h.Hostname != "" {
		p.Hostname = &h.Hostname
	}
	if h.SNMPVersion != "" {
		p.SNMPVersion = &h.SNMPVersion
	}
	if h.SNMPCommunity != "" {
		p.SNMPCommunity = &h.SNMPCommunity
	}
	if h.SNMPPort > 0 {
		p.SNMPPort = &h.SNMPPort
	}
	if h.SNMPTimeout > 0 {
		p.SNMPTimeout = &h.SNMPTimeout
	}
	if h.SSHPort > 0 {
		p.SSHPort = &h.SSHPort
	}
	if h.SSHUsername != "" {
		p.SSHUsername = &h.SSHUsername
	}

	const q = `
		INSERT INTO monitored_hosts
			(id, tenant_id, name, display_name, ip_address, hostname, monitoring_type,
			 snmp_version, snmp_community, snmp_port, snmp_timeout_secs,
			 ssh_port, ssh_username, status)
		VALUES
			(:id, :tenant_id, :name, :display_name, :ip_address, :hostname, :monitoring_type,
			 :snmp_version, :snmp_community, :snmp_port, :snmp_timeout_secs,
			 :ssh_port, :ssh_username, :status)
		RETURNING created_at, updated_at`

	stmt, args, err := r.db.BindNamed(q, p)
	if err != nil {
		return fmt.Errorf("bind named: %w", err)
	}

	row := struct {
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}{}
	if err := r.db.QueryRowxContext(ctx, stmt, args...).StructScan(&row); err != nil {
		return fmt.Errorf("insert monitored_host: %w", err)
	}
	h.CreatedAt = row.CreatedAt
	h.UpdatedAt = row.UpdatedAt
	return nil
}

func (r *HostRepository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.MonitoredHost, error) {
	var row monitoredHostRow
	q := `SELECT ` + hostCols + ` FROM monitored_hosts WHERE tenant_id=$1 AND id=$2 AND deleted_at IS NULL`
	if err := r.db.GetContext(ctx, &row, q, tenantID, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}
	return row.toDomain(), nil
}

func (r *HostRepository) GetByAssetID(ctx context.Context, tenantID, assetID uuid.UUID) (*domain.MonitoredHost, error) {
	var row monitoredHostRow
	q := `SELECT ` + hostCols + ` FROM monitored_hosts WHERE tenant_id=$1 AND asset_id=$2 AND deleted_at IS NULL LIMIT 1`
	if err := r.db.GetContext(ctx, &row, q, tenantID, assetID); err != nil {
		if err == sql.ErrNoRows {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}
	return row.toDomain(), nil
}

func (r *HostRepository) Update(ctx context.Context, h *domain.MonitoredHost) error {
	const q = `
		UPDATE monitored_hosts SET
			name=$3, display_name=$4, ip_address=$5, hostname=$6,
			monitoring_type=$7, status=$8, updated_at=NOW()
		WHERE tenant_id=$1 AND id=$2 AND deleted_at IS NULL`

	var dispName, ipAddr, host *string
	if h.DisplayName != "" {
		dispName = &h.DisplayName
	}
	if h.IPAddress != "" {
		ipAddr = &h.IPAddress
	}
	if h.Hostname != "" {
		host = &h.Hostname
	}
	_, err := r.db.ExecContext(ctx, q, h.TenantID, h.ID, h.Name, dispName, ipAddr, host, string(h.MonitoringType), string(h.Status))
	return err
}

func (r *HostRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.HostCheckStatus, checkedAt time.Time) error {
	const q = `UPDATE monitored_hosts SET last_status=$2, last_check_at=$3, updated_at=NOW() WHERE id=$1 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, q, id, string(status), checkedAt)
	return err
}

func (r *HostRepository) SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error {
	const q = `UPDATE monitored_hosts SET deleted_at=NOW() WHERE tenant_id=$1 AND id=$2 AND deleted_at IS NULL`
	res, err := r.db.ExecContext(ctx, q, tenantID, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return shared.ErrNotFound
	}
	return nil
}

func (r *HostRepository) List(ctx context.Context, tenantID uuid.UUID, filter domain.MonitoredHostFilter, page shared.Page) (shared.PageResult[*domain.MonitoredHost], error) {
	page = page.Normalized()

	where := "WHERE tenant_id=$1 AND deleted_at IS NULL"
	args := []any{tenantID}
	idx := 2

	if filter.Status != nil {
		where += fmt.Sprintf(" AND status=$%d", idx)
		args = append(args, string(*filter.Status))
		idx++
	}
	if filter.LastStatus != nil {
		where += fmt.Sprintf(" AND last_status=$%d", idx)
		args = append(args, string(*filter.LastStatus))
		idx++
	}
	if filter.Search != "" {
		where += fmt.Sprintf(" AND (name ILIKE $%d OR ip_address::text ILIKE $%d)", idx, idx+1)
		args = append(args, "%"+filter.Search+"%", "%"+filter.Search+"%")
		idx += 2
	}

	var total int64
	if err := r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM monitored_hosts `+where, args...); err != nil {
		return shared.PageResult[*domain.MonitoredHost]{}, err
	}

	listQ := fmt.Sprintf(
		`SELECT %s FROM monitored_hosts %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		hostCols, where, idx, idx+1,
	)
	args = append(args, page.Size, page.Offset())

	var rows []monitoredHostRow
	if err := r.db.SelectContext(ctx, &rows, listQ, args...); err != nil {
		return shared.PageResult[*domain.MonitoredHost]{}, err
	}

	items := make([]*domain.MonitoredHost, len(rows))
	for i := range rows {
		items[i] = rows[i].toDomain()
	}
	return shared.NewPageResult(items, total, page), nil
}

func (r *HostRepository) GetDownHosts(ctx context.Context, tenantID uuid.UUID) ([]*domain.MonitoredHost, error) {
	q := `SELECT ` + hostCols + ` FROM monitored_hosts WHERE tenant_id=$1 AND last_status='down' AND deleted_at IS NULL ORDER BY last_check_at ASC LIMIT 50`
	var rows []monitoredHostRow
	if err := r.db.SelectContext(ctx, &rows, q, tenantID); err != nil {
		return nil, err
	}
	hosts := make([]*domain.MonitoredHost, len(rows))
	for i := range rows {
		hosts[i] = rows[i].toDomain()
	}
	return hosts, nil
}

func (r *HostRepository) CountByStatus(ctx context.Context, tenantID uuid.UUID) (map[domain.HostCheckStatus]int64, error) {
	type countRow struct {
		Status string `db:"last_status"`
		Count  int64  `db:"cnt"`
	}
	const q = `
		SELECT COALESCE(last_status, 'unknown') AS last_status, COUNT(*) AS cnt
		FROM monitored_hosts WHERE tenant_id=$1 AND deleted_at IS NULL
		GROUP BY last_status`

	var rows []countRow
	if err := r.db.SelectContext(ctx, &rows, q, tenantID); err != nil {
		return nil, err
	}
	out := make(map[domain.HostCheckStatus]int64)
	for _, row := range rows {
		out[domain.HostCheckStatus(row.Status)] = row.Count
	}
	return out, nil
}
