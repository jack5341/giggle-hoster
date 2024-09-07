package cf

import (
	"context"
	"errors"
	"os"

	"github.com/cloudflare/cloudflare-go"
)

type Cloudflare struct {
	Client *cloudflare.API
}

var (
	ErrCloudflareConnectionCouldNotBeEstablished = errors.New("cloudflare connection could not be established please check your CLOUDFLARE_API_KEY and CLOUDFLARE_API_EMAIL variables")
	ErrDNSRecordCouldNotBeCreated                = errors.New("DNS record could not be created")
	ErrDNSRecordCouldNotBeDeleted                = errors.New("DNS record could not be deleted")
)

func (c *Cloudflare) NewClient() (*Cloudflare, error) {
	api, err := cloudflare.New(os.Getenv("CLOUDFLARE_API_KEY"), os.Getenv("CLOUDFLARE_API_EMAIL"))
	if err != nil {
		return nil, errors.Join(ErrCloudflareConnectionCouldNotBeEstablished, err)
	}

	c.Client = api

	return c, nil
}

func (c *Cloudflare) CreateRecord(ctx context.Context, record *cloudflare.ResourceContainer, params cloudflare.CreateDNSRecordParams) (cloudflare.DNSRecord, error) {
	client := c.Client
	dnsRecord, err := client.CreateDNSRecord(ctx, record, params)
	if err != nil {
		return cloudflare.DNSRecord{}, errors.Join(ErrDNSRecordCouldNotBeCreated, err)
	}

	return dnsRecord, nil
}

func (c *Cloudflare) DeleteRecord(ctx context.Context, record *cloudflare.ResourceContainer, recordID string) error {
	client := c.Client
	if err := client.DeleteDNSRecord(ctx, record, recordID); err != nil {
		return errors.Join(ErrDNSRecordCouldNotBeDeleted, err)
	}

	return nil
}
