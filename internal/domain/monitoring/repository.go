package monitoring

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/infracore/infracore/internal/domain/shared"
)

type MonitoredHostFilter struct {
	SiteID         *uuid.UUID
	AgentID        *uuid.UUID
	MonitoringType *MonitoringType
	Status         *HostStatus
	LastStatus     *HostCheckStatus
	Search         string
}

type MonitoredHostRepository interface {
	Create(ctx context.Context, h *MonitoredHost) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*MonitoredHost, error)
	GetByAssetID(ctx context.Context, tenantID, assetID uuid.UUID) (*MonitoredHost, error)
	Update(ctx context.Context, h *MonitoredHost) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status HostCheckStatus, checkedAt time.Time) error
	SoftDelete(ctx context.Context, tenantID, id uuid.UUID) error
	List(ctx context.Context, tenantID uuid.UUID, filter MonitoredHostFilter, page shared.Page) (shared.PageResult[*MonitoredHost], error)
	GetDownHosts(ctx context.Context, tenantID uuid.UUID) ([]*MonitoredHost, error)
	CountByStatus(ctx context.Context, tenantID uuid.UUID) (map[HostCheckStatus]int64, error)
}

type CheckDefinitionRepository interface {
	Create(ctx context.Context, c *CheckDefinition) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*CheckDefinition, error)
	Update(ctx context.Context, c *CheckDefinition) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
	ListByHost(ctx context.Context, tenantID, hostID uuid.UUID) ([]*CheckDefinition, error)
	GetDueChecks(ctx context.Context, tenantID uuid.UUID, before time.Time) ([]*CheckDefinition, error)
}

// TimeRange specifies a window for metric queries.
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// AggregationType defines how to aggregate metric data points.
type AggregationType string

const (
	AggregationAvg AggregationType = "avg"
	AggregationMax AggregationType = "max"
	AggregationMin AggregationType = "min"
	AggregationSum AggregationType = "sum"
	AggregationLast AggregationType = "last"
)

// MetricQueryResult holds a single aggregated data point.
type MetricQueryResult struct {
	Time  time.Time
	Value float64
}

type MetricRepository interface {
	// Write a batch of metric points (optimized for high throughput).
	WriteBatch(ctx context.Context, points []MetricPoint) error
	// Query aggregated metrics over a time range with bucketing.
	Query(ctx context.Context, tenantID, hostID uuid.UUID, metricName string,
		tr TimeRange, bucketDuration time.Duration, agg AggregationType) ([]MetricQueryResult, error)
	// Get the latest single value for a host/metric combination.
	GetLatest(ctx context.Context, tenantID, hostID uuid.UUID, metricName string) (*MetricPoint, error)
	// Get distinct metric names available for a host.
	GetMetricNames(ctx context.Context, tenantID, hostID uuid.UUID) ([]string, error)
}

type AvailabilityRepository interface {
	Create(ctx context.Context, r *AvailabilityRecord) error
	CloseRecord(ctx context.Context, hostID uuid.UUID, endedAt time.Time) error
	GetCurrent(ctx context.Context, hostID uuid.UUID) (*AvailabilityRecord, error)
	GetHistory(ctx context.Context, tenantID, hostID uuid.UUID, tr TimeRange, page shared.Page) (shared.PageResult[*AvailabilityRecord], error)
	GetUptimePercent(ctx context.Context, tenantID, hostID uuid.UUID, tr TimeRange) (float64, error)
}

type ServiceCheckRepository interface {
	Create(ctx context.Context, s *ServiceCheck) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*ServiceCheck, error)
	Update(ctx context.Context, s *ServiceCheck) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
	ListByHost(ctx context.Context, tenantID, hostID uuid.UUID) ([]*ServiceCheck, error)
}
