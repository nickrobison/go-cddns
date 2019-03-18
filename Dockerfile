FROM alpine:latest
LABEL version="0.1.0"
LABEL description="Docker image for running go-cddns application"
LABEL maintainer="Nick Robison nick@nickrobison.com"

# Link MUSL in place of GLIBC
# https://stackoverflow.com/questions/34729748/installed-go-binary-not-found-in-path-on-alpine-linux-docker
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

# Setup the necessary SSL certs
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true

# Add the go binary and config files
COPY bin/go-cddns*_amd64 /app/go-cddns
WORKDIR /app
CMD ["./go-cddns", "--config=/etc/config.json"]
