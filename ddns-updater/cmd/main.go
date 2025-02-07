package main

import (
	"os"
	"fmt"
	"log"
	"time"
	"strconv"
	
	"ddns-updater/internal/aws"
	"ddns-updater/internal/network"
)

func main() {
	pingInterval, err := strconv.ParseInt(os.Getenv("PING_INTERVAL"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	hostedZoneID := os.Getenv("hostedZoneID")
	recordName := os.Getenv("recordName")

	dnsIp, err := aws.GetARecord(hostedZoneID, recordName)
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
			log.Printf("[DDNS-UPDATER] Received public IP %s\n", publicIp)
			if dnsIp != publicIp {
				log.Print("[DDNS-UPDATER] A records and public IP don't match, updating...")
				err = aws.UpdateARecord(hostedZoneID, recordName, publicIp)
				if err != nil {
					fmt.Print(err)
				} else {
					log.Printf("[DDNS-UPDATER] Successfully updated A records in Route 53 to %s\n", dnsIp)
				}
				dnsIp = publicIp
			}
		}

		time.Sleep(time.Duration(pingInterval) * time.Second)
	}
}
