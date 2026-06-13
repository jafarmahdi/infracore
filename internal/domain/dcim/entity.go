package dcim

import (
	"time"

	"github.com/google/uuid"
)

// ─────────────────────────────────────────────────────────────
// DataCenter  — aggregate root
// ─────────────────────────────────────────────────────────────

type DataCenter struct {
	ID                  uuid.UUID
	TenantID            uuid.UUID
	SiteID              uuid.UUID
	Name                string
	Description         string
	FacilityCode        string
	PhysicalAddress     string
	NOCContact          string
	PowerCapacityKVA    float64
	CoolingCapacityTons float64
	Status              DataCenterStatus
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time
}

type DataCenterStatus string

const (
	DataCenterStatusActive         DataCenterStatus = "active"
	DataCenterStatusPlanned        DataCenterStatus = "planned"
	DataCenterStatusDecommissioned DataCenterStatus = "decommissioned"
)

// ─────────────────────────────────────────────────────────────
// Room
// ─────────────────────────────────────────────────────────────

type Room struct {
	ID           uuid.UUID
	TenantID     uuid.UUID
	DataCenterID uuid.UUID
	Name         string
	Description  string
	Floor        string
	RoomNumber   string
	Dimensions   *RoomDimensions
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

type RoomDimensions struct {
	WidthMeters  float64 `json:"width_meters"`
	LengthMeters float64 `json:"length_meters"`
	HeightMeters float64 `json:"height_meters"`
}

func (r *RoomDimensions) AreaSqMeters() float64 {
	return r.WidthMeters * r.LengthMeters
}

// ─────────────────────────────────────────────────────────────
// Rack  — aggregate root (owns rack units)
// ─────────────────────────────────────────────────────────────

type Rack struct {
	ID            uuid.UUID
	TenantID      uuid.UUID
	SiteID        uuid.UUID
	RoomID        *uuid.UUID
	Name          string
	Description   string
	FacilityID    string
	RackType      RackType
	Status        RackStatus
	WidthInches   int
	TotalUnits    int
	UnitNumbering RackUnitNumbering
	Position      *RackPosition
	MaxWeightKg   *float64
	MaxPowerWatts *int
	Comment       string
	Tags          []string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

type RackType string

const (
	RackTypeOpenFrame RackType = "open-frame"
	RackTypeEnclosed  RackType = "enclosed"
	RackTypeWallMount RackType = "wall-mount"
	RackType2Post     RackType = "2-post"
	RackType4Post     RackType = "4-post"
)

type RackStatus string

const (
	RackStatusActive         RackStatus = "active"
	RackStatusPlanned        RackStatus = "planned"
	RackStatusReserved       RackStatus = "reserved"
	RackStatusDecommissioned RackStatus = "decommissioned"
)

type RackUnitNumbering string

const (
	RackUnitBottomToTop RackUnitNumbering = "bottom-to-top"
	RackUnitTopToBottom RackUnitNumbering = "top-to-bottom"
)

type RackPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// ─────────────────────────────────────────────────────────────
// PowerFeed
// ─────────────────────────────────────────────────────────────

type PowerFeed struct {
	ID                 uuid.UUID
	TenantID           uuid.UUID
	RackID             uuid.UUID
	Name               string
	Supply             PowerSupplyType
	Phase              PowerPhase
	Voltage            int
	Amperage           int
	MaxUtilizationPct  int
	Status             PowerFeedStatus
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type PowerSupplyType string

const (
	PowerSupplyAC PowerSupplyType = "AC"
	PowerSupplyDC PowerSupplyType = "DC"
)

type PowerPhase string

const (
	PowerPhaseSingle PowerPhase = "single-phase"
	PowerPhaseThree  PowerPhase = "three-phase"
)

type PowerFeedStatus string

const (
	PowerFeedStatusActive  PowerFeedStatus = "active"
	PowerFeedStatusOffline PowerFeedStatus = "offline"
	PowerFeedStatusPlanned PowerFeedStatus = "planned"
	PowerFeedStatusFailed  PowerFeedStatus = "failed"
)

// TotalWatts returns the feed's rated power in watts.
func (p *PowerFeed) TotalWatts() int {
	return p.Voltage * p.Amperage
}

// MaxAllocatedWatts returns the derated safe maximum.
func (p *PowerFeed) MaxAllocatedWatts() int {
	return p.TotalWatts() * p.MaxUtilizationPct / 100
}

// ─────────────────────────────────────────────────────────────
// PDU
// ─────────────────────────────────────────────────────────────

type PDU struct {
	ID            uuid.UUID
	TenantID      uuid.UUID
	RackID        uuid.UUID
	PowerFeedID   *uuid.UUID
	Name          string
	Model         string
	Manufacturer  string
	SerialNumber  string
	PDUType       PDUType
	TotalOutlets  int
	Amperage      int
	Voltage       int
	Status        string
	IPAddress     string
	ManagementURL string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

type PDUType string

const (
	PDUTypeBasic     PDUType = "basic"
	PDUTypeMetered   PDUType = "metered"
	PDUTypeMonitored PDUType = "monitored"
	PDUTypeSwitched  PDUType = "switched"
	PDUTypeSmart     PDUType = "smart"
)

// ─────────────────────────────────────────────────────────────
// PatchPanel
// ─────────────────────────────────────────────────────────────

type PatchPanel struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	RackID         uuid.UUID
	Name           string
	PortCount      int
	PortType       string
	RackUnitStart  int
	RackUnitHeight int
	Label          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

// ─────────────────────────────────────────────────────────────
// Cable
// ─────────────────────────────────────────────────────────────

type Cable struct {
	ID           uuid.UUID
	TenantID     uuid.UUID
	CableType    CableType
	Status       CableStatus
	Label        string
	Color        string
	LengthMeters *float64
	Description  string
	TerminationA CableTermination
	TerminationB CableTermination
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CableType string

const (
	CableTypeCat5e  CableType = "cat5e"
	CableTypeCat6   CableType = "cat6"
	CableTypeCat6a  CableType = "cat6a"
	CableTypeCat7   CableType = "cat7"
	CableTypeCat8   CableType = "cat8"
	CableTypeMMF    CableType = "mmf"
	CableTypeSMF    CableType = "smf"
	CableTypeCoax   CableType = "coax"
	CableTypePower  CableType = "power"
	CableTypeDAC    CableType = "dac"
	CableTypeAOC    CableType = "aoc"
)

type CableStatus string

const (
	CableStatusConnected      CableStatus = "connected"
	CableStatusPlanned        CableStatus = "planned"
	CableStatusDecommissioned CableStatus = "decommissioned"
)

// CableTermination identifies one end of a cable connection.
// It is polymorphic: TerminationType indicates which domain table holds the termination.
type CableTermination struct {
	TerminationType string    `json:"termination_type"` // "network_interface" | "patch_panel_port" | "pdu_outlet"
	TerminationID   uuid.UUID `json:"termination_id"`
}
