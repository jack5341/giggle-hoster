package types

import (
	"time"

	"github.com/google/uuid"
)

type Provider string
type Status string

const (
	STARTING Status = "STARTING"
	PENDING  Status = "PENDING"
	WELL     Status = "WELL"
	STOPPED  Status = "STOPPED"
)

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
	Provider  Provider  `gorm:"notnull;Provider"`
	Size      string    `gorm:"notnull"`
	Url       string    `gorm:"notnull"`
	Status    Status    `gorm:"notnull;Status"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
