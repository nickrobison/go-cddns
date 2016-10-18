# go-cddns
Golang client for dynamically updating cloudflare DNS records on a specified interval. Useful if you're using Cloudflare to point to a device with a dynamic IP Address

## Usage

Create a config.json with the following structure:

```json
{
  "UpdateInterval": "{interval (in minutes) to check for an updated IP Address}",
  "Key": "{Cloudflare API Key}",
  "Email": "{Cloudflare Email Address}",
  "DomainName": "{Cloudflare domain to modify}",
  "RecordName": "{Array of DNS records to update}"
  }
  ```
