package shared

import (
	"time"

	"github.com/google/uuid"
)

// DomainEvent is the base interface for all domain events.
type DomainEvent interface {
	EventName() string
	OccurredAt() time.Time
	AggregateID() uuid.UUID
	TenantID() uuid.UUID
}

// BaseEvent provides common domain event fields.
type BaseEvent struct {
	name        string
	aggregateID uuid.UUID
	tenantID    uuid.UUID
	occurredAt  time.Time
}

func NewBaseEvent(name string, aggregateID, tenantID uuid.UUID) BaseEvent {
	return BaseEvent{
		name:        name,
		aggregateID: aggregateID,
		tenantID:    tenantID,
		occurredAt:  time.Now().UTC(),
	}
}

func (e BaseEvent) EventName() string       { return e.name }
func (e BaseEvent) OccurredAt() time.Time   { return e.occurredAt }
func (e BaseEvent) AggregateID() uuid.UUID  { return e.aggregateID }
func (e BaseEvent) TenantID() uuid.UUID     { return e.tenantID }

// EventBus is the interface for publishing and subscribing to domain events.
type EventBus interface {
	Publish(event DomainEvent) error
	Subscribe(eventName string, handler EventHandler)
}

// EventHandler processes a domain event.
type EventHandler func(event DomainEvent) error
