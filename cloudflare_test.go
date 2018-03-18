package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateDNSRecord(t *testing.T) {
	server := httptest.NewServer((http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {

		defer r.Body.Close()
		// Get the JSON as a map
		var updateRecord map[string]*json.RawMessage
		json.NewDecoder(r.Body).Decode(&updateRecord)

		fmt.Println(updateRecord)

		var proxied bool
		json.Unmarshal(*updateRecord["proxied"], &proxied)
		if proxied != true {
			t.Error("Proxied should be true, and is not")
		}

		var ttl int
		json.Unmarshal(*updateRecord["ttl"], &ttl)
		if ttl != 130 {
			t.Error("TTL value should match")
		}

		var content string
		json.Unmarshal(*updateRecord["content"], &content)
		if content != "127.0.0.11" {
			t.Error("Record not updated correctly")
		}
		w.Write([]byte("{\"success\": true, \"result\": {\"name\": \"test\", \"type\": \"A\"}}"))
	})))
	defer server.Close()

	testConfig := &config{
		CloudflareURL: server.URL,
		ZoneID:        "test-zone",
	}

	cfm := createManager(testConfig)
	testDNSRecord := cloudflareDNSRecord{
		Name:    "dns-one",
		Proxied: true,
		Content: "127.0.0.1",
		Ttl:     130,
	}
	cfm.updateDNSRecord(testDNSRecord, "127.0.0.11")
}

func createManager(userConfig *config) *cloudflareManager {
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
		*userConfig,
		client,
		"",
	}
	return cfm
}
