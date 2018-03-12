VERSION := 0.0.1
PKGNAME := go-cddns
LICENSE := MIT
URL := http://github.com/nickrobison/go-cddns
RELEASE := 4
USER := cddns
DESC := Dynamic DNS client for Cloudflare
MAINTAINER := Nick Robison <nick@nickrobison.com>
DOCKER_WDIR := /tmp/fpm
DOCKER_FPM := fpm-ubuntu
PLATFORMS := linux/amd64 linux/arm linux/arm64

FPM_OPTS=-s dir -v $(VERSION) -n $(PKGNAME) \
  --license "$(LICENSE)" \
  --vendor "$(VENDOR)" \
  --maintainer "$(MAINTAINER)" \
  --url "$(URL)" \
  --description  "$(DESC)" \
  --verbose
#	--config-files etc/ \
  --verbose

DEB_OPTS= -t deb --deb-user $(USER) --deb-no-default-config-files

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

build-all: $(PLATFORMS)

$(PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) go build -o 'bin/go-cddns_0.0.1_$(arch)' .
.PHONY: build-all $(PLATFORMS)

.PHONY: release
release: 
	# Build
	GOOS=linux GOARCH=amd64 go build .
	# Package
	docker run --rm -it -v "${PWD}:${DOCKER_WDIR}" -w ${DOCKER_WDIR} fpm-ubuntu ${DEB_OPTS} \
	--iteration ${RELEASE} \
	--architecture amd64 \
	${FPM_OPTS} \
	go-cddns
	# Remove it
	rm go-cddns
	# Upload it
	# Manually, for now
	# curl -H "X-Bintray-Debian-Distribution: jessie,xenial,stretch" \
	# -H "X-Bintray-Debian-Component: main" \
	# -H "X-Bintray-Debian-Architecture: amd64" \
	# -unickrobison:${API_KEY} -T go-cddns_${VERSION}-${RELEASE}_amd64.deb \
	# https://api.bintray.com/content/nickrobison/debian/go-cddns/0.0.1/go-cddns/go-cddns_${VERSION}-${RELEASE}_amd64.deb;deb_distribution=jessie,xenial;deb_component=main;deb_architecture=amd64
	# Build
	GOOS=linux GOARCH=arm go build .
	# Package
	docker run --rm -it -v "${PWD}:${DOCKER_WDIR}" -w ${DOCKER_WDIR} fpm-ubuntu ${DEB_OPTS} \
	--iteration ${RELEASE} \
	--architecture armhf \
	${FPM_OPTS} \
	go-cddns
	# Remove it
	rm go-cddns
	# Upload everything
	# Manually, for now
	# curl -H "X-Bintray-Debian-Distribution: jessie,xenial,stretch" \
	# -H "X-Bintray-Debian-Component: main" \
	# -H "X-Bintray-Debian-Architecture: amd64" \
	# -unickrobison:${API_KEY} -T go-cddns_${VERSION}-${RELEASE}_armhf.deb \
	# https://api.bintray.com/content/nickrobison/debian/go-cddns/0.0.1/go-cddns/go-cddns_${VERSION}-${RELEASE}_armhf.deb
