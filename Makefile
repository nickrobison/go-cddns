VERSION := 0.1.0
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

build: test
	go build -o 'go-cddns' .

test:
	go test -v ./...

all: $(PLATFORMS)

$(PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) go build -o 'bin/go-cddns_${VERSION}_$(arch)' .

clean:
	-rm -rf bin/
	-rm go-cddns
	-rm *.deb
	-rm packaging/debian/usr/bin/go-cddns

bintray:
	docker pull alanfranz/fpm-within-docker:debian-jessie
	# Copy the amd64 binary
	cp bin/go-cddns_${VERSION}_amd64 #{PWD}packaging/debian/usr/bin/go-cddns
	# Package amd64
	docker run --rm -it -v "${PWD}:${DOCKER_WDIR}" -w ${DOCKER_WDIR} --entrypoint fpm alanfranz/fpm-within-docker:debian-jessie ${DEB_OPTS} \
	--iteration ${RELEASE} \
	--architecture amd64 \
	--deb-systemd go-cddns.service \
	-C packaging/debian \
	${FPM_OPTS}
	# Upload it
	./upload.sh ${VERSION} ${RELEASE} amd64
	# Copy the arm binary
	cp bin/go-cddns_${VERSION}_arm ${PWD}/packaging/debian/usr/bin/go-cddns
	# Package arm
	docker run --rm -it -v "${PWD}:${DOCKER_WDIR}" -w ${DOCKER_WDIR} --entrypoint fpm alanfranz/fpm-within-docker:debian-jessie ${DEB_OPTS} \
	--iteration ${RELEASE} \
	--architecture armhf \
	--deb-systemd go-cddns.service \
	-C packaging/debian \
	${FPM_OPTS}
	# Upload everything
	./upload.sh ${VERSION} ${RELEASE} armhf

docker:
	docker build . -t nickrobison/go-cddns
	./docker_push.sh

release: all bintray docker


.PHONY: all $(PLATFORMS) release clean test bintray docker
