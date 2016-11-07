package main

import (
	"net/http"
	"time"
)

type cddnsContext struct {
	HttpClient       *http.Client
	Config           config
	CurrentIPAddress string
}

type cloudflareMetadata struct {
	auto_added bool
}

type cloudflareDNSRecord struct {
	Id          string
	RecordType  string `json:"type"`
	Name        string
	Content     string
	Proxiable   bool
	Proxied     bool
	Ttl         int
	Locked      bool
	Zone_id     string
	Zone_name   string
	Created_on  time.Time
	Modified_on time.Time
	Data        interface{}
	Meta        cloudflareMetadata
	Priority    int
}

type cloudflarePlan struct {
	Id            string
	Name          string
	Price         int
	Currency      string
	Frequency     string
	Legacy_id     string
	Is_subscribed bool
	Can_subscribe bool
}

type host struct {
	Name    string
	Website string
}

type cloudflareZone struct {
	Id                    string
	Name                  string
	Development_mode      int
	Original_name_servers []string
	Original_registrar    string
	Original_dnshost      string
	Created_on            time.Time
	Modified_on           time.Time
	Name_servers          []string
	Owner                 interface{}
	Permissions           []string
	Plan                  cloudflarePlan
	Plan_pending          cloudflarePlan
	Status                string
	Paused                bool
	Type                  string
	Checked_on            time.Time
	Host                  host
	Vanity_name_servers   []string
	Betas                 []string
	Deactivation_reason   string
	Meta                  interface{}
}

type cloudflareResultInfo struct {
	Page        int
	Per_page    int
	Count       int
	Total_count int
}

type cloudflareResponse struct {
	Success     bool
	Errors      []string
	Messages    []string
	Result_info cloudflareResultInfo
}

type cloudflareZoneResponse struct {
	cloudflareResponse
	Result []cloudflareZone
}

type cloudflareDNSResponse struct {
	cloudflareResponse
	Result []cloudflareDNSRecord
}

type cloudflareDNSUpdateResponse struct {
	cloudflareResponse
	Result map[string]interface{}
}

type newCloudflareDNSRecord struct {
	Type    dnsType `json:"type"`
	Name    string  `json:"name"`
	Content string  `json:"content"`
}

type config struct {
	UpdateInterval int
	Key            string
	Email          string
	DomainName     string
	ZoneID         string
	RecordNames    []string
	Remove         bool
}
