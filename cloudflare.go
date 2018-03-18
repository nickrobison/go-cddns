package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type httpMethod string
type dnsType string

const (
	// 	HTTP Methods
	GET    httpMethod = "GET"
	PUT    httpMethod = "PUT"
	POST   httpMethod = "POST"
	DELETE httpMethod = "DELETE"
	//	DNS Types
	A     dnsType = "A"
	AAAA  dnsType = "AAAA"
	CNAME dnsType = "CNAME"
	TXT   dnsType = "TXT"
	SRV   dnsType = "SRV"
	LOC   dnsType = "LOC"
	MX    dnsType = "MX"
	NS    dnsType = "NS"
	SPF   dnsType = "SPF"
)

type cloudflareManager struct {
	config
	Client           *http.Client
	CurrentIPAddress string
}

func (c *cloudflareManager) GetDNSRecords() (map[string]cloudflareDNSRecord, error) {
	//	Create the URL
	var cloudflareURL *url.URL
	var err error
	cloudflareURL, err = cloudflareURL.Parse(c.CloudflareURL)
	if err != nil {
		logger.Fatalln(err)
	}
	req, err := c.buildHTTPRequest(cloudflareURL, GET, nil)
	if err != nil {
		return nil, err
	}

	//Add the search params
	params := req.URL.Query()
	params.Add("name", c.DomainName)
	req.URL.RawQuery = params.Encode()

	zoneResponse, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer zoneResponse.Body.Close()

	var zones cloudflareZoneResponse
	err = json.NewDecoder(zoneResponse.Body).Decode(&zones)
	if err != nil {
		return nil, err
	}

	//	Create the URL
	cloudflareURL, err = cloudflareURL.Parse(c.CloudflareURL + zones.Result[0].Id + "/dns_records")
	if err != nil {
		return nil, err
	}
	req, err = c.buildHTTPRequest(cloudflareURL, GET, nil)
	if err != nil {
		return nil, err
	}
	dnsResponse, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer dnsResponse.Body.Close()
	var records cloudflareDNSResponse
	err = json.NewDecoder(dnsResponse.Body).Decode(&records)
	if err != nil {
		return nil, err
	}

	var dnsMap = make(map[string]cloudflareDNSRecord)

	for _, record := range records.Result {
		dnsMap[record.Name] = record
	}

	return dnsMap, nil
}

func (c *cloudflareManager) updateDNSRecord(record cloudflareDNSRecord, ipAddress string) {
	// Update the record with the new ip address
	record.Content = ipAddress

	//	Create the URL
	var cloudflareURL *url.URL
	var err error
	cloudflareURL, err = cloudflareURL.Parse(c.CloudflareURL + c.ZoneID + "/dns_records/" + record.Id)
	if err != nil {
		logger.Fatalln(err)
	}

	newRecord := &newCloudflareDNSRecord{
		Type:    A,
		Name:    record.Name,
		Content: ipAddress,
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(newRecord)
	req, err := c.buildHTTPRequest(cloudflareURL, PUT, b)
	if err != nil {
		logger.Fatalln(err)
	}
	dnsResponse, err := c.Client.Do(req)
	if err != nil {
		logger.Fatalln(err)
	}
	defer dnsResponse.Body.Close()
	var updatedRecord cloudflareDNSUpdateResponse
	err = json.NewDecoder(dnsResponse.Body).Decode(&updatedRecord)
	if err != nil {
		logger.Fatalln(err)
	}

	if !updatedRecord.Success {
		logger.Fatalf("Unable to updated record %v. %v", record.Name, updatedRecord.Messages)
	}

}

func (c *cloudflareManager) CreateDNSRecord(recordName string, dns dnsType, ipAddress string) {
	record := &newCloudflareDNSRecord{
		Type:    dns,
		Name:    recordName,
		Content: ipAddress,
	}

	//	Create the URL
	var cloudflareURL *url.URL
	var err error
	cloudflareURL, err = cloudflareURL.Parse(c.CloudflareURL + c.ZoneID + "/dns_records")
	if err != nil {
		logger.Fatalln(err)
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(record)

	var req *http.Request
	req, err = c.buildHTTPRequest(cloudflareURL, POST, b)
	if err != nil {
		logger.Fatalln(err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		logger.Fatalln(err)
	}
	defer resp.Body.Close()

	var respContent map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&respContent)
	if err != nil {
		logger.Fatalln(err)
	}

	if !respContent["success"].(bool) {
		logger.Fatalf("Unable to create record %v. %v\n", recordName, respContent["messages"])
	}
}

func (c *cloudflareManager) DeleteDNSRecord(record cloudflareDNSRecord) {
	//	Create the URL
	var cloudflareURL *url.URL
	var err error
	cloudflareURL, err = cloudflareURL.Parse(c.CloudflareURL + c.ZoneID + "/dns_records/" + record.Id)
	if err != nil {
		logger.Fatalln(err)
	}

	req, err := c.buildHTTPRequest(cloudflareURL, DELETE, nil)
	if err != nil {
		logger.Fatalln(err)
	}
	dnsResponse, err := c.Client.Do(req)
	if err != nil {
		logger.Fatalln(err)
	}
	defer dnsResponse.Body.Close()
	var deletedRecord cloudflareResponse
	err = json.NewDecoder(dnsResponse.Body).Decode(&deletedRecord)
	if err != nil {
		logger.Fatalln(err)
	}

	if !deletedRecord.Success {
		logger.Fatalf("Unable to delete record %v. %v", record.Name, deletedRecord.Messages)
	}
}

func (c *cloudflareManager) buildHTTPRequest(url *url.URL, method httpMethod, body io.Reader) (*http.Request, error) {
	var req *http.Request
	var err error
	switch method {
	case GET:
		{
			req, err = http.NewRequest(http.MethodGet, url.String(), nil)
			if err != nil {
				return nil, err
			}
		}
	case POST:
		{
			req, err = http.NewRequest(http.MethodPost, url.String(), body)
			if err != nil {
				return nil, err
			}

		}
	case PUT:
		{
			req, err = http.NewRequest(http.MethodPut, url.String(), body)
			if err != nil {
				return nil, err
			}
		}
	case DELETE:
		{
			req, err = http.NewRequest(http.MethodDelete, url.String(), nil)
			if err != nil {
				return nil, err
			}
		}
	}
	//	Set the headers
	req.Header.Set("X-Auth-Email", c.Email)
	req.Header.Set("X-Auth-Key", c.Key)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
