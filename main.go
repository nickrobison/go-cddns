package main

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var logger *log.Logger

func main() {
	logger = log.New(os.Stdout, "go-cddns:", log.Lshortfile)

	//	Read config file
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		logger.Fatalln("Cannot open config.json")
	}

	var cloudflareConfig config
	err = json.Unmarshal(file, &cloudflareConfig)
	if err != nil {
		logger.Fatalln(err)
	}

	logger.Printf("Updating records for: %v", cloudflareConfig.RecordNames)

	//	Check for correct config values
	if cloudflareConfig.UpdateInterval < 5 {
		logger.Fatalln("Update interval cannot be less than 5 minutes")
	}
	logger.Printf("Setting update interval to %v minutes\n", cloudflareConfig.UpdateInterval)

	minimalTLSConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	tlsTransport := &http.Transport{
		TLSClientConfig: minimalTLSConfig,
	}
	client := &http.Client{
		Transport: tlsTransport,
	}

	cfm := &cloudflareManager{
		Client:           client,
		Config:           cloudflareConfig,
		CurrentIPAddress: "",
	}

	updateCloudflareRecord(cfm)
	for _ = range time.NewTicker(time.Duration(cloudflareConfig.UpdateInterval) * time.Minute).C {
		updateCloudflareRecord(cfm)
	}

	logger.Println("done")
}

func updateCloudflareRecord(cfm *cloudflareManager) {

	// Get an updated record
	updatedIpAddress, err := getCurrentIPAddress(cfm.Client)
	if err != nil {
		logger.Println(err)
	}

	if updatedIpAddress != cfm.CurrentIPAddress {
		logger.Printf("IP Address changed from %v, to %v\n", cfm.CurrentIPAddress, updatedIpAddress)
		cfm.CurrentIPAddress = updatedIpAddress

		//	Get the cloudflare DNS records
		records, err := cfm.GetDNSRecords()
		if err != nil {
			logger.Println(err)
		}
		//	Range over the given RecordNames and update them, if possible
		for _, name := range cfm.Config.RecordNames {
			if record, ok := records[name]; ok {
				if record.Content != cfm.CurrentIPAddress {
					logger.Printf("Record %v has ip address %v, updating to %v\n", name, record.Content, cfm.CurrentIPAddress)
					cfm.updateDNSRecord(record, cfm.CurrentIPAddress)
				} else {
					logger.Printf("Record %v has correct IP Address, skipping\n", name)
				}
			} else {
				//	Record doesn't exist, create it
				logger.Printf("Record %v doesn't exist, creating it with IP Address: %v\n", name, cfm.CurrentIPAddress)
				cfm.CreateDNSRecord(name, A, cfm.CurrentIPAddress)
			}
		}
	}
}

func getCurrentIPAddress(httpClient *http.Client) (string, error) {
	resp, err := httpClient.Get("http://bot.whatismyipaddress.com")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ipAddress, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(ipAddress), nil

}
