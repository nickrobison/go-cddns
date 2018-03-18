VERSION := 0.0.1
PKGNAME := go-cddns
LICENSE := MIT
URL := http://github.com/nickrobison/go-cddns
RELEASE := 2
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

DEB_OPTS= -t deb --deb-user $(USER) --after-install packaging/debian/go-cddns.postinst

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
	GOOS=linux GOARCH=amd64 go build -o packaging/debian/usr/bin/go-cddns .
	# Package
	docker run --rm -it -v "${PWD}:${DOCKER_WDIR}" -w ${DOCKER_WDIR} --entrypoint fpm alanfranz/fpm-within-docker:debian-jessie ${DEB_OPTS} \
	--iteration ${RELEASE} \
	--architecture amd64 \
	--deb-systemd go-cddns.service \
	-C packaging/debian \
	${FPM_OPTS} \
	# Remove it
	rm packaging/debian/usr/bin/go-cddns
	# Upload it
	./upload.sh ${VERSION} ${RELEASE} amd64
	# Build
	GOOS=linux GOARCH=arm go build -o packaging/debian/usr/bin/go-cddns .
	# Package
	docker run --rm -it -v "${PWD}:${DOCKER_WDIR}" -w ${DOCKER_WDIR} --entrypoint fpm alanfranz/fpm-within-docker:debian-jessie ${DEB_OPTS} \
	--iteration ${RELEASE} \
	--architecture armhf \
	--deb-systemd go-cddns.service \
	-C packaging/debian \
	${FPM_OPTS} \
	# Remove it
	rm packaging/debian/usr/bin/go-cddns
	# Upload everything
	./upload.sh ${VERSION} ${RELEASE} armhf	
