[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.org/nickrobison/go-cddns.svg?branch=master)](https://travis-ci.org/nickrobison/go-cddns)
[ ![Download](https://api.bintray.com/packages/nickrobison/debian/go-cddns/images/download.svg) ](https://bintray.com/nickrobison/debian/go-cddns/_latestVersion)
# go-cddns
Golang client for dynamically updating cloudflare DNS records on a specified interval. Useful if you're using Cloudflare to point to a device with a dynamic IP Address

## Installation
go get -u github.com/nickrobison/go-cddns

### Debian repository
We also now have a debian (and Ubuntu) repo with builds for both amd64 and armhf architectures.
Since we require systemd, the builds only support debian jessie and newer, and ubuntu xenial (16.04) and later.

The repository is hosted on bintray, so there are some special setup instructions.

```bash
sudo apt-get install apt-transport-https # Bintray only supports https connections
echo "deb https://dl.bintray.com/nickrobison/debian {distribution} {components}" | sudo tee -a /etc/apt/sources.list
apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 379CE192D401AB61 # We need to import the Bintray public key
sudo apt-get update && apt-get install go-cddns
```

The package installation will create a config file in ```/etc/go-cddns/config.json```, which is where you should set your configuration options.

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

Run the application, optionally specifying the path to the config file.

```bash
go-cddns -file=/path/to/file
```

## Notes

* The update interval must be more than 5 minutes, per the WhatIsMyIP API [rules](http://whatismyipaddress.com/api).
* The records names must be FQDNs, even though they don't appear in the cloudflare dashboard as such.
* If the Remove field is set to true, the listed DNS records will be removed when the program exits.