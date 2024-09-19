package cf

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type Cloudflare struct {
	APIURL   string
	APIKey   string
	APIToken string
}

var (
	ErrCloudflareConnectionCouldNotBeEstablished = errors.New("cloudflare connection could not be established. Please check your API token and endpoint.")
	ErrDNSRecordCouldNotBeCreated                = errors.New("DNS record could not be created")
	ErrDNSRecordCouldNotBeDeleted                = errors.New("DNS record could not be deleted")
	ErrEnvVariablesHaveToBeConfigured            = errors.New("API token environment variable must be configured")
)

func NewClient() (*Cloudflare, error) {
	apiToken := os.Getenv("CLOUDFLARE_API_TOKEN")

	if apiToken == "" {
		return nil, ErrEnvVariablesHaveToBeConfigured
	}

	return &Cloudflare{
		APIURL:   "https://api.cloudflare.com/client/v4",
		APIToken: apiToken,
	}, nil
}

func (c *Cloudflare) CreateRecord(ctx context.Context, zoneID string, params any) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/zones/%s/dns_records", c.APIURL, zoneID)
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Join(ErrCloudflareConnectionCouldNotBeEstablished, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result["success"] != true {
		return nil, ErrDNSRecordCouldNotBeCreated
	}

	return result, nil
}

func (c *Cloudflare) DeleteRecord(ctx context.Context, zoneID, recordID string) error {
	url := fmt.Sprintf("%s/zones/%s/dns_records/%s", c.APIURL, zoneID, recordID)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Join(ErrCloudflareConnectionCouldNotBeEstablished, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
