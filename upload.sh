#!/bin/bash

set -e 

curl -H "X-Bintray-Debian-Distribution: jessie,xenial,stretch" \
	-H "X-Bintray-Debian-Component: main" \
	-H "X-Bintray-Debian-Architecture: ${3}" \
	-unickrobison:${API_KEY} -T go-cddns_${1}-${2}_${3}.deb \
	https://api.bintray.com/content/nickrobison/debian/go-cddns/${1}/go-cddns_${1}-${2}_${3}.deb;publish=0