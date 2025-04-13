package main

import (
	"os"
	"fmt"
	"log"
	"time"
	"strconv"

	"ddns-updater/internal/network"
	"ddns-updater/internal/cloudflare"
)

func main() {
	pingInterval, err := strconv.ParseInt(os.Getenv("PING_INTERVAL"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	apiToken := os.Getenv("CF_DNS_API_TOKEN")
	zoneId := os.Getenv("CLOUDFLARE_ZONE_ID")
	recordName := os.Getenv("CLOUDFLARE_RECORD_NAME")

	dnsIp, err := cloudflare.GetARecord(apiToken, zoneId, recordName)
	if err != nil {
		log.Fatalf("[DDNS-UPDATER] Error while querying public IP %s", err)
	} else {
		log.Printf("[DDNS-UPDATER] DNS provider A record found : %s\n", dnsIp)
	}

	for {
		publicIp, err := network.GetPublicIp()
		if err != nil {
			log.Printf("[DDNS-UPDATER] Error retrieving public IP: %v", err)
		} else {
			log.Printf("[DDNS-UPDATER] Public IP is: %s\n", publicIp)
			if dnsIp != publicIp {
				log.Print("[DDNS-UPDATER] A records and public IP don't match, updating...")
				err = cloudflare.UpdateARecord(apiToken, zoneId, recordName, publicIp)
				if err != nil {
					fmt.Print(err)
				} else {
					log.Printf("[DDNS-UPDATER] Successfully updated A records in Route 53 to %s\n", publicIp)
				}
				dnsIp = publicIp
			}
		}

		time.Sleep(time.Duration(pingInterval) * time.Second)
	}
}
