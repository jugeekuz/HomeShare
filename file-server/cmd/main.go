// main.go
package main

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"time"

	"file-server/config"
	"file-server/internal/app"
	"file-server/internal/helpers"
	"file-server/internal/job"
)

type AcmeJSON struct {
    LetsEncrypt struct {
        Account struct {
            Email         string `json:"Email"`
            Registration  struct {
                Body struct {
                    Status  string   `json:"status"`
                    Contact []string `json:"contact"`
                } `json:"body"`
                Uri string `json:"uri"`
            } `json:"Registration"`
            PrivateKey string `json:"PrivateKey"`
            KeyType    string `json:"KeyType"`
        } `json:"Account"`
        Certificates []struct {
            Domain struct {
                Main string `json:"main"`
            } `json:"domain"`
            Certificate string `json:"certificate"`
            Key         string `json:"key"`
            Store       string `json:"Store"`
        } `json:"Certificates"`
    } `json:"letencrypt"`
}
func loadCertificate(domain string) (tls.Certificate, error) {
    data, err := os.ReadFile("certs/acme.json")
    if err != nil {
        return tls.Certificate{}, err
    }

    var acme AcmeJSON
    if err := json.Unmarshal(data, &acme); err != nil {
        return tls.Certificate{}, err
    }

    for _, certEntry := range acme.LetsEncrypt.Certificates {
        if certEntry.Domain.Main == domain {
			// Decode Base64-encoded certificate and key
			certPEM, err := base64.StdEncoding.DecodeString(certEntry.Certificate)
			if err != nil {
				return tls.Certificate{}, fmt.Errorf("failed to decode certificate: %w", err)
			}
			keyPEM, err := base64.StdEncoding.DecodeString(certEntry.Key)
			if err != nil {
				return tls.Certificate{}, fmt.Errorf("failed to decode key: %w", err)
			}
			return tls.X509KeyPair(certPEM, keyPEM)
		}
    }

    return tls.Certificate{}, fmt.Errorf("no certificate found for domain %s", domain)
}

func main() {
	cfg := config.LoadConfig()
	if err := os.MkdirAll(cfg.SharingDir, os.ModePerm); err != nil {
		log.Fatal("[FILE-SERVER] Error while creating sharing folder directory")
	}

	go func () {
		for {
			_ = helpers.CleanupExpiredFolders(cfg.SharingDir)
			time.Sleep(30 * time.Minute)
		}
	}()
	
	job_timeout := 45 * time.Second
	jm := job.NewJobManager(job_timeout)

	server, err := app.SetupServer(jm, app.InitDatabase)
	if err != nil {
		log.Fatalf("[FILE-SERVER] Server setup failed: %v", err)
	}

	server.Addr = ":443"

	cert, err := loadCertificate(cfg.Domain)
	if err != nil {
		log.Fatalf("Failed to load certificate: %v", err)
	}
	
	tlsConfig := &tls.Config{
        Certificates: []tls.Certificate{cert},
    }

	server.TLSConfig = tlsConfig

	log.Fatal(server.ListenAndServeTLS("","",))
}
