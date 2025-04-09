package cloudflare

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go"
)

func GetARecord(apiToken string, zoneID string, recordName string) (string, error) {
	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		return "", fmt.Errorf("failed to create Cloudflare API client: %v", err)
	}

	filter := cloudflare.ListDNSRecordsParams{
		Type: "A",
		Name: recordName,
	}

	records, _, err := api.ListDNSRecords(context.Background(), cloudflare.ZoneIdentifier(zoneID), filter)
	if err != nil {
		return "", fmt.Errorf("failed to list DNS records: %v", err)
	}

	if len(records) == 0 {
		return "", fmt.Errorf("no A record found for %s", recordName)
	}

	return records[0].Content, nil
}

func UpdateARecord(apiToken string, zoneID string, recordName string, newIp string) error {
	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		return fmt.Errorf("failed to create Cloudflare API client: %v", err)
	}

	filter := cloudflare.ListDNSRecordsParams{
		Type: "A",
		Name: recordName,
	}

	records, _, err := api.ListDNSRecords(context.Background(), cloudflare.ZoneIdentifier(zoneID), filter)
	if err != nil {
		return fmt.Errorf("failed to list DNS records: %v", err)
	}

	if len(records) == 0 {
		return fmt.Errorf("no A record found for %s", recordName)
	}

	recordID := records[0].ID

	dnsRecord := cloudflare.UpdateDNSRecordParams{
		ID:			recordID,
		Type:    	"A",
		Name:   	recordName,
		Content: 	newIp,
		TTL:     	300,
		Proxied: 	records[0].Proxied,
	}

	if _, err := api.UpdateDNSRecord(context.Background(), cloudflare.ZoneIdentifier(zoneID), dnsRecord); err != nil {
		return fmt.Errorf("failed to update A record: %v", err)
	}

	return nil
}
