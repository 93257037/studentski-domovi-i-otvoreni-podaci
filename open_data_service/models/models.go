package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StDom represents a student dormitory (read-only from st_dom_service)
type StDom struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Ime             string             `bson:"ime" json:"ime"`
	Address         string             `bson:"address" json:"address"`
	TelephoneNumber string             `bson:"telephone_number" json:"telephone_number"`
	Email           string             `bson:"email" json:"email"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

// Soba represents a room (read-only from st_dom_service)
type Soba struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	StDomID    primitive.ObjectID `bson:"st_dom_id" json:"st_dom_id"`
	Krevetnost int                `bson:"krevetnost" json:"krevetnost"`
	Luksuzi    []string           `bson:"luksuzi" json:"luksuzi"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

// Aplikacija represents a room application (read-only)
type Aplikacija struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	BrojIndexa string             `bson:"broj_indexa" json:"broj_indexa"`
	Prosek     int                `bson:"prosek" json:"prosek"`
	SobaID     primitive.ObjectID `bson:"soba_id" json:"soba_id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	IsActive   bool               `bson:"is_active" json:"is_active"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

// PrihvacenaAplikacija represents an accepted application (read-only)
type PrihvacenaAplikacija struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AplikacijaID primitive.ObjectID `bson:"aplikacija_id" json:"aplikacija_id"`
	UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`
	BrojIndexa   string             `bson:"broj_indexa" json:"broj_indexa"`
	Prosek       int                `bson:"prosek" json:"prosek"`
	SobaID       primitive.ObjectID `bson:"soba_id" json:"soba_id"`
	AcademicYear string             `bson:"academic_year" json:"academic_year"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

// ====================
// Open Data Response Models
// ====================

// PublicStatistics - Aggregated public statistics about all dorms
type PublicStatistics struct {
	TotalDorms              int                        `json:"total_dorms"`
	TotalRooms              int                        `json:"total_rooms"`
	TotalCapacity           int                        `json:"total_capacity"`
	TotalOccupied           int                        `json:"total_occupied"`
	OccupancyRate           float64                    `json:"occupancy_rate"`
	AmenitiesDistribution   map[string]int             `json:"amenities_distribution"`
	RoomTypeDistribution    map[int]int                `json:"room_type_distribution"`
	ApplicationStatistics   ApplicationStats           `json:"application_statistics"`
	DormStatistics          []DormStats                `json:"dorm_statistics"`
	LastUpdated             time.Time                  `json:"last_updated"`
}

// ApplicationStats - Statistics about applications
type ApplicationStats struct {
	TotalApplications         int     `json:"total_applications"`
	ActiveApplications        int     `json:"active_applications"`
	AcceptedApplications      int     `json:"accepted_applications"`
	AcceptanceRate            float64 `json:"acceptance_rate"`
	AverageGradeOfAccepted    float64 `json:"average_grade_of_accepted"`
	AverageGradeOfApplications float64 `json:"average_grade_of_applications"`
}

// DormStats - Statistics for a specific dorm
type DormStats struct {
	DormID              primitive.ObjectID `json:"dorm_id"`
	DormName            string             `json:"dorm_name"`
	Address             string             `json:"address"`
	TotalRooms          int                `json:"total_rooms"`
	TotalCapacity       int                `json:"total_capacity"`
	OccupiedSpots       int                `json:"occupied_spots"`
	AvailableSpots      int                `json:"available_spots"`
	OccupancyRate       float64            `json:"occupancy_rate"`
	RoomTypes           map[int]int        `json:"room_types"`
	Amenities           map[string]int     `json:"amenities"`
}

// RoomAvailability - Public room availability with filters
type RoomAvailability struct {
	RoomID         primitive.ObjectID `json:"room_id"`
	DormID         primitive.ObjectID `json:"dorm_id"`
	DormName       string             `json:"dorm_name"`
	DormAddress    string             `json:"dorm_address"`
	Capacity       int                `json:"capacity"`
	Occupied       int                `json:"occupied"`
	AvailableSpots int                `json:"available_spots"`
	Amenities      []string           `json:"amenities"`
	IsAvailable    bool               `json:"is_available"`
}

// DormComparison - Compare multiple dorms side-by-side
type DormComparison struct {
	Dorms           []DormComparisonDetail `json:"dorms"`
	ComparisonDate  time.Time              `json:"comparison_date"`
}

// DormComparisonDetail - Detailed comparison data for a dorm
type DormComparisonDetail struct {
	DormID              primitive.ObjectID  `json:"dorm_id"`
	DormName            string              `json:"dorm_name"`
	Address             string              `json:"address"`
	ContactInfo         ContactInfo         `json:"contact_info"`
	Capacity            CapacityInfo        `json:"capacity"`
	RoomDistribution    map[int]int         `json:"room_distribution"`
	AmenitiesOffered    map[string]int      `json:"amenities_offered"`
	ApplicationMetrics  ApplicationMetrics  `json:"application_metrics"`
}

// ContactInfo - Contact information for a dorm
type ContactInfo struct {
	Phone string `json:"phone"`
	Email string `json:"email"`
}

// CapacityInfo - Capacity information
type CapacityInfo struct {
	TotalRooms      int     `json:"total_rooms"`
	TotalCapacity   int     `json:"total_capacity"`
	OccupiedSpots   int     `json:"occupied_spots"`
	AvailableSpots  int     `json:"available_spots"`
	OccupancyRate   float64 `json:"occupancy_rate"`
}

// ApplicationMetrics - Application metrics for a dorm
type ApplicationMetrics struct {
	TotalApplications    int     `json:"total_applications"`
	AcceptedApplications int     `json:"accepted_applications"`
	AcceptanceRate       float64 `json:"acceptance_rate"`
	AverageGrade         float64 `json:"average_grade"`
}

// ApplicationTrends - Historical trends of applications
type ApplicationTrends struct {
	YearlyTrends    []YearlyTrend         `json:"yearly_trends"`
	DormTrends      []DormApplicationTrend `json:"dorm_trends"`
	OverallMetrics  TrendMetrics          `json:"overall_metrics"`
	GeneratedAt     time.Time             `json:"generated_at"`
}

// YearlyTrend - Application trend for a specific academic year
type YearlyTrend struct {
	AcademicYear         string  `json:"academic_year"`
	TotalApplications    int     `json:"total_applications"`
	AcceptedApplications int     `json:"accepted_applications"`
	AcceptanceRate       float64 `json:"acceptance_rate"`
	AverageGrade         float64 `json:"average_grade"`
	MinGrade             int     `json:"min_grade"`
	MaxGrade             int     `json:"max_grade"`
}

// DormApplicationTrend - Application trend for a specific dorm
type DormApplicationTrend struct {
	DormID               primitive.ObjectID `json:"dorm_id"`
	DormName             string             `json:"dorm_name"`
	TotalApplications    int                `json:"total_applications"`
	AcceptedApplications int                `json:"accepted_applications"`
	AcceptanceRate       float64            `json:"acceptance_rate"`
}

// TrendMetrics - Overall trend metrics
type TrendMetrics struct {
	TotalYears              int     `json:"total_years"`
	AverageApplicationsPerYear int  `json:"average_applications_per_year"`
	TrendDirection          string  `json:"trend_direction"` // "increasing", "decreasing", "stable"
}

// OccupancyHeatmap - Real-time occupancy visualization data
type OccupancyHeatmap struct {
	Dorms       []DormOccupancyPoint `json:"dorms"`
	Summary     OccupancySummary     `json:"summary"`
	GeneratedAt time.Time            `json:"generated_at"`
}

// DormOccupancyPoint - Occupancy data point for a dorm
type DormOccupancyPoint struct {
	DormID          primitive.ObjectID `json:"dorm_id"`
	DormName        string             `json:"dorm_name"`
	Address         string             `json:"address"`
	OccupancyRate   float64            `json:"occupancy_rate"`
	Status          string             `json:"status"` // "high", "medium", "low"
	TotalCapacity   int                `json:"total_capacity"`
	OccupiedSpots   int                `json:"occupied_spots"`
	AvailableSpots  int                `json:"available_spots"`
}

// OccupancySummary - Summary of occupancy across all dorms
type OccupancySummary struct {
	AverageOccupancy float64 `json:"average_occupancy"`
	HighestOccupancy float64 `json:"highest_occupancy"`
	LowestOccupancy  float64 `json:"lowest_occupancy"`
	FullDorms        int     `json:"full_dorms"`
	EmptyDorms       int     `json:"empty_dorms"`
}

// ExportFormat - Format for data export
type ExportFormat string

const (
	ExportFormatJSON ExportFormat = "json"
	ExportFormatCSV  ExportFormat = "csv"
)

// ExportDataRequest - Request for data export
type ExportDataRequest struct {
	Dataset string       `form:"dataset" binding:"required"` // "dorms", "rooms", "statistics"
	Format  ExportFormat `form:"format" binding:"required"`
}

// RoomSearchFilters - Filters for room search
type RoomSearchFilters struct {
	DormID         string   `form:"dorm_id"`
	MinCapacity    int      `form:"min_capacity"`
	MaxCapacity    int      `form:"max_capacity"`
	Amenities      []string `form:"amenities"`
	OnlyAvailable  bool     `form:"only_available"`
	Limit          int      `form:"limit"`
	Offset         int      `form:"offset"`
}

// DormSearchFilters - Filters for dorm search
type DormSearchFilters struct {
	MinOccupancy float64 `form:"min_occupancy"`
	MaxOccupancy float64 `form:"max_occupancy"`
	HasAvailability bool `form:"has_availability"`
	Limit        int     `form:"limit"`
	Offset       int     `form:"offset"`
}
