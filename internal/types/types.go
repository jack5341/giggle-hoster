package types

import (
	"time"

	"github.com/google/uuid"
)

type Provider string

const (
	HCLOUD Provider = "HCLOUD"
	// TODO: Implement new cloud providers
	// AWS    Provider = "AWS"
	// GCLOUD Provider = "GCLOUD"
)

type HcloudGetAllServiceType struct {
	ServerTypes []HcloudServerType `json:"server_types"`
}

type HcloudServerType struct {
	Architecture    string            `json:"architecture"`
	Cores           int               `json:"cores"`
	CPUType         string            `json:"cpu_type"`
	Deprecated      bool              `json:"deprecated"`
	Deprecation     HcloudDeprecation `json:"deprecation"`
	Description     string            `json:"description"`
	Disk            int               `json:"disk"`
	ID              int               `json:"id"`
	IncludedTraffic interface{}       `json:"included_traffic"`
	Memory          int               `json:"memory"`
	Name            string            `json:"name"`
	Prices          []HcloudPrice     `json:"prices"`
	StorageType     string            `json:"storage_type"`
}

type HcloudDeprecation struct {
	Announced        time.Time `json:"announced"`
	UnavailableAfter time.Time `json:"unavailable_after"`
}

type HcloudPrice struct {
	IncludedTraffic   int                `json:"included_traffic"`
	Location          string             `json:"location"`
	PriceHourly       HcloudPriceDetails `json:"price_hourly"`
	PriceMonthly      HcloudPriceDetails `json:"price_monthly"`
	PricePerTbTraffic HcloudPriceDetails `json:"price_per_tb_traffic"`
}

type HcloudPriceDetails struct {
	Gross string `json:"gross"`
	Net   string `json:"net"`
}

type Node struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name      string    `gorm:"notnull"`
	Metadata  []string  `gorm:"type:text[]"`
	Provider  Provider  `gorm:"notnull"`
	Size      string    `gorm:"notnull"`
	Pods      []Pod     `gorm:"foreignKey:BelongsToNodeID"`
	FreeMem   int       `gorm:"notnull"`
	FreeCPU   int       `gorm:"notnull"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Pod struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name            string    `gorm:"notnull"`
	Metadata        []string  `gorm:"type:text[]"`
	ExposedPort     int       `gorm:"notnull"`
	RequestedMem    int       `gorm:"notnull"`
	RequestedCPU    int       `gorm:"notnull"`
	BelongsToUser   uuid.UUID `gorm:"notnull"`
	BelongsToNodeID uuid.UUID `gorm:"notnull"`
}
