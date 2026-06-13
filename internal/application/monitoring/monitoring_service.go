package monitoring

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	domain "github.com/infracore/infracore/internal/domain/monitoring"
	"github.com/infracore/infracore/internal/domain/shared"
	"go.uber.org/zap"
)

type MonitoringService struct {
	hostRepo domain.MonitoredHostRepository
	log      *zap.Logger
}

func NewMonitoringService(hostRepo domain.MonitoredHostRepository, log *zap.Logger) *MonitoringService {
	return &MonitoringService{hostRepo: hostRepo, log: log}
}

// ── DTOs ─────────────────────────────────────────────────────────────────────

type CreateHostRequest struct {
	Name           string `json:"name"            binding:"required,min=1,max=255"`
	DisplayName    string `json:"display_name"`
	IPAddress      string `json:"ip_address"`
	Hostname       string `json:"hostname"`
	MonitoringType string `json:"monitoring_type"  binding:"required,oneof=agent snmp wmi icmp ssh api"`
	SNMPVersion    string `json:"snmp_version"`
	SNMPCommunity  string `json:"snmp_community"`
	SNMPPort       int    `json:"snmp_port"`
	SNMPTimeout    int    `json:"snmp_timeout_secs"`
	SSHPort        int    `json:"ssh_port"`
	SSHUsername    string `json:"ssh_username"`
}

type HostResponse struct {
	ID             uuid.UUID  `json:"id"`
	TenantID       uuid.UUID  `json:"tenant_id"`
	Name           string     `json:"name"`
	DisplayName    string     `json:"display_name"`
	IPAddress      string     `json:"ip_address"`
	Hostname       string     `json:"hostname"`
	MonitoringType string     `json:"monitoring_type"`
	Status         string     `json:"status"`
	LastStatus     *string    `json:"last_status"`
	UptimePercent  *float64   `json:"uptime_percent"`
	LastCheckAt    *time.Time `json:"last_check_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type ListHostsRequest struct {
	Search     string `form:"search"`
	Status     string `form:"status"`
	LastStatus string `form:"last_status"`
	Page       int    `form:"page"`
	PageSize   int    `form:"page_size"`
}

type ListHostsResponse struct {
	Items      []HostResponse `json:"items"`
	TotalItems int64          `json:"total_items"`
	TotalPages int            `json:"total_pages"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
}

func toHostResponse(h *domain.MonitoredHost) HostResponse {
	r := HostResponse{
		ID:             h.ID,
		TenantID:       h.TenantID,
		Name:           h.Name,
		DisplayName:    h.DisplayName,
		IPAddress:      h.IPAddress,
		Hostname:       h.Hostname,
		MonitoringType: string(h.MonitoringType),
		Status:         string(h.Status),
		UptimePercent:  h.UptimePercent,
		LastCheckAt:    h.LastCheckAt,
		CreatedAt:      h.CreatedAt,
		UpdatedAt:      h.UpdatedAt,
	}
	if h.LastStatus != nil {
		s := string(*h.LastStatus)
		r.LastStatus = &s
	}
	return r
}

// ── Service methods ───────────────────────────────────────────────────────────

func (s *MonitoringService) CreateHost(ctx context.Context, tenantID uuid.UUID, req CreateHostRequest) (HostResponse, error) {
	mt := domain.MonitoringType(strings.ToLower(req.MonitoringType))

	h := &domain.MonitoredHost{
		TenantID:       tenantID,
		Name:           strings.TrimSpace(req.Name),
		DisplayName:    strings.TrimSpace(req.DisplayName),
		IPAddress:      strings.TrimSpace(req.IPAddress),
		Hostname:       strings.TrimSpace(req.Hostname),
		MonitoringType: mt,
		SNMPVersion:    req.SNMPVersion,
		SNMPCommunity:  req.SNMPCommunity,
		SNMPPort:       req.SNMPPort,
		SNMPTimeout:    req.SNMPTimeout,
		SSHPort:        req.SSHPort,
		SSHUsername:    req.SSHUsername,
		Status:         domain.HostStatusPending,
	}

	if err := s.hostRepo.Create(ctx, h); err != nil {
		s.log.Error("create monitored host", zap.Error(err))
		return HostResponse{}, err
	}
	return toHostResponse(h), nil
}

func (s *MonitoringService) GetHost(ctx context.Context, tenantID, id uuid.UUID) (HostResponse, error) {
	h, err := s.hostRepo.GetByID(ctx, tenantID, id)
	if err != nil {
		return HostResponse{}, err
	}
	return toHostResponse(h), nil
}

func (s *MonitoringService) ListHosts(ctx context.Context, tenantID uuid.UUID, req ListHostsRequest) (ListHostsResponse, error) {
	filter := domain.MonitoredHostFilter{Search: req.Search}
	if req.Status != "" {
		st := domain.HostStatus(req.Status)
		filter.Status = &st
	}
	if req.LastStatus != "" {
		ls := domain.HostCheckStatus(req.LastStatus)
		filter.LastStatus = &ls
	}

	page := shared.NewPage(req.Page, req.PageSize)
	result, err := s.hostRepo.List(ctx, tenantID, filter, page)
	if err != nil {
		return ListHostsResponse{}, err
	}

	items := make([]HostResponse, len(result.Items))
	for i, h := range result.Items {
		items[i] = toHostResponse(h)
	}
	return ListHostsResponse{
		Items:      items,
		TotalItems: result.TotalItems,
		TotalPages: result.TotalPages,
		Page:       result.Page,
		PageSize:   result.PageSize,
	}, nil
}

func (s *MonitoringService) DeleteHost(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.hostRepo.SoftDelete(ctx, tenantID, id)
}

func (s *MonitoringService) GetStatusCounts(ctx context.Context, tenantID uuid.UUID) (map[string]int64, error) {
	counts, err := s.hostRepo.CountByStatus(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	out := make(map[string]int64)
	for k, v := range counts {
		out[string(k)] = v
	}
	return out, nil
}
