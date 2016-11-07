[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
# go-cddns
Golang client for dynamically updating cloudflare DNS records on a specified interval. Useful if you're using Cloudflare to point to a device with a dynamic IP Address

## Installation
go get -u github.com/nickrobison/go-cddns

## Usage

Create a config.json with the following structure:

```json
{
  "UpdateInterval": "{interval (in minutes) to check for an updated IP Address}",
  "Key": "{Cloudflare API Key}",
  "Email": "{Cloudflare Email Address}",
  "DomainName": "{Cloudflare domain to modify}",
  "RecordName": "{Array of DNS records to update}",
  "Remove": "{Boolean of whether or not to remove the records on shutdown}"
  }
  ```
## Notes

* The update interval must be more than 5 minutes, per the WhatIsMyIP API [rules](http://whatismyipaddress.com/api).
* The records names must be FQDNs, even though they don't appear in the cloudflare dashboard as such.
* If the Remove field is set to true, the listed DNS records will be removed when the program exits.

