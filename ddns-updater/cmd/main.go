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
		fmt.Print(err)
		fmt.Print("error while receiving")
	} else {
		log.Printf("Starting ip found : %s\n", dnsIp)
	}

	for {
		publicIp, err := network.GetPublicIp()
		if err != nil {
			log.Printf("Error retrieving IP: %v", err)
		} else {
			if dnsIp != publicIp {
				err = aws.UpdateARecord(hostedZoneID, recordName, publicIp)
				if err != nil {
					fmt.Print(err)
				} else {
					log.Printf("Updated A records in Route 53 to %s\n", dnsIp)
				}
				dnsIp = publicIp
			}
			log.Printf("Received IP %s\n", publicIp)
		}

		time.Sleep(time.Duration(pingInterval) * time.Second)
	}
}
